Ich möchte, dass du einen Crossplane Provider (https://www.crossplane.io/) für die Artefakte innerhalb von dem Azure Data Explorer baust. Also für sowas wie eine ADX-Tabelle, eine Material SU Update Policy, eine Function, aber auch sowas wie Row Level Security, Retention Policy, Caching Policy, External Table, Continuous Export, also alles, was sonst so drumherum auch an Artefakten anfällt. Ich werde dir da einen Terraform-Provider bereitstellen.

Ein Link dazu. Da kannst du dich dann orientieren, aber ich möchte halt, wie gesagt, dass es halt ein Crossplane Provider wird und dass du den komplett ausschmückst und ausbaust. 

Terraform Provider der ein paar Artefakte schon kann: https://github.com/favoretti/terraform-provider-azure-adx

Wenn dir noch weitere Artefakte oder Ressourcen auf dem ARDX-Cluster einfallen, dann kannst du das gerne innerhalb von diesem Crossplane-Provider implementieren.  

---

# Kritische Bewertung & offene Fragen (Stand 2026-09-07)

Lesart des Textes oben: "Material SU Update Policy" = Materialized View + Update Policy, "ARDX" = ADX.

## 1. Ist das machbar?

**Ja.** Alle ADX-Artefakte sind über Management-Commands (`.create` / `.alter` / `.show` / `.drop`) am Kusto-Endpoint steuerbar. Das ist genau das Muster von provider-sql (Observe per `SHOW`, Create/Update per DDL). Es gibt aber vier Punkte, die das Projekt tragen oder kippen. Sie sind nicht "ob", sondern "wie" – und sie müssen vor dem ersten Controller entschieden sein.

### Fakten, die die Planung ändern

- **Der TF-Provider ist nur als Kommando-Referenz brauchbar, nicht als Code-Vorlage.** favoretti/terraform-provider-azure-adx nutzt `azure-kusto-go v0.7.0` und `go-autorest` (seit 2023 deprecated). Das aktuelle Kusto-SDK ist `github.com/Azure/azure-kusto-go/azkustodata` (v1, komplett anderes Package-Layout, Auth über `azidentity`). Übernehmen können wir: welche KQL-Commands pro Artefakt, wie `.show`-Ausgaben aussehen, wo Drift-Vergleich schwierig ist. Nicht übernehmen: SDK-Layer, Auth, Terraform-Plugin-Logik.
- **`provider-azure-adx-old` ist wertlos als Basis.** Es ist ein unverändertes provider-template (Initial commit, April 2025) auf crossplane-runtime **v1.16**. Ordner `tablesecurityrole/` und `examples/adx/` sind leer. Das aktuelle provider-template steht auf `crossplane-runtime/v2 v2.4.0`, Go 1.25, controller-runtime 0.23. → Frisch vom aktuellen Template starten, alten Ordner löschen.
- **Crossplane ist bei v2.4.** Cluster-scoped MRs sind der Legacy-Pfad und werden abgekündigt. Ein neues Projekt sollte namespaced MRs liefern (siehe F9).
- **Es gibt einen Kusto-Emulator** (`mcr.microsoft.com/azuredataexplorer/kustainer-linux`, Testcontainers-Modul vorhanden), aber: keine Authentifizierung, keine Streaming Ingestion, External Tables nur auf lokale Dateien, Retention/Partitioning werden akzeptiert aber nicht angewendet, kein Python-Plugin. Der Free Cluster kann keine External Tables. Ein Teil der Artefakte ist also **nur gegen einen bezahlten Cluster** testbar (Security Roles, External Tables auf Storage, Continuous Export, Managed-Identity-Policy, Streaming-Ingestion-Policy in Wirkung).

### Die vier harten Probleme

**P1 – Drift-Erkennung von KQL-Text.** Function-Body, Update-Policy-Query, RLS-Query, MV-Query, Continuous-Export-Query sind Freitext. Wenn `.show` etwas anderes zurückgibt als im YAML steht (Zeilenenden, Einrückung, Parameter-Formatierung `(a:string,b:int)` vs. `(a:string, b:int)`), läuft der Controller in eine Endlosschleife aus `.create-or-alter`. Braucht eine definierte Normalisierung (mindestens: Trim, CRLF→LF, Parameter-Liste kanonisieren) und Tests dafür. Aggressives Whitespace-Kollabieren ist gefährlich (Whitespace in String-Literalen ist semantisch). Signal, dass es schiefgeht: MR flappt zwischen `Synced=True/False`, hohe Command-Rate im Cluster.

**P2 – Destruktive Änderungen ohne "Plan".** Terraform zeigt einen Plan, den ein Mensch absegnet. Crossplane reconciliert sofort. Bei ADX heißt das:
- `.alter table` mit weniger Spalten **löscht Daten** dieser Spalten. `.alter-merge` kann nur hinzufügen.
- Spaltentyp-Änderung ist nicht sauber möglich (`.alter column` macht bestehende Daten unlesbar).
- Änderung der MV-Query ist nur eingeschränkt per `.alter materialized-view` möglich, sonst Drop+Recreate mit Backfill.
- Umbenennen = Drop + Create aus Crossplane-Sicht.
- `.drop table` löscht alle Daten.
Guardrails müssen ins API (z. B. `schemaUpdateMode: Merge|Replace`, Default `Merge`; Typänderung → Fehler im Status statt Auto-Drop). Ohne das ist der Provider in Produktion nicht vertretbar.

**P3 – Grenze zu ARM.** Cluster, Database, Data Connections (Event Hub/IoT Hub/Event Grid/Cosmos), Principal Assignments, Database Scripts, Attached Databases, Managed Private Endpoints sind ARM-Ressourcen und existieren bereits in provider-upbound-azure (`kusto.azure.upbound.io`). Dieser Provider sollte **nur** das abdecken, was hinter dem Kusto-Endpoint liegt. Nebeneffekt: Microsofts eigene "deklarative" Lösung (`kusto_script`, ARM Database Script) ist Fire-and-forget ohne Observe/Drift – genau die Lücke, die wir füllen. Das gehört als Begründung in die README.

**P4 – Skalierung / Throttling.** Jede MR pollt per `.show`. 1.000 MRs mit Crossplane-Default-Poll von 1 Minute = 1.000 Management-Commands/Minute pro Cluster. Kusto hat Concurrency-Limits für Management-Commands (Capacity Policy, Request Rate Limit Policy). Braucht: Poll-Intervall-Default 10 Minuten (wie Upbound), begrenzte `MaxConcurrentReconciles`, saubere Behandlung von HTTP 429 / "throttled" (Backoff, nicht als Fehler zählen). Signal: `.show commands` zeigt Spitzen, Cluster-Admins beschweren sich.

### Weitere Stolpersteine (kleiner, aber real)

- **KQL-Injection.** Namen und Bodies werden in Commands interpoliert. Namen konsequent als `['name']` quoten, Policies als JSON serialisieren statt String-Konkatenation, Query-Bodies nur in `<|`-Form bzw. `{ }`-Blöcken.
- **Secrets in External Tables.** Connection Strings mit SAS/Account Key. Kusto verschleiert nur Strings mit `h@'...'`-Prefix in `.show`. Folge: Wenn wir verschleiern (müssen wir), ist auf dem Connection String **kein Drift-Vergleich** möglich. Klartext-Secrets dürfen nie in `status.atProvider` landen. Primärer Pfad sollte Managed Identity (`;managed_identity=system`) sein, Secret-Auth als Write-only-Feld.
- **Principal-Normalisierung.** `.show ... principals` liefert `aadapp=<objectId>;<tenantId>`, der User schreibt aber `aadapp=<appId>;<tenant>` oder `aaduser=name@domain`. Naiver Vergleich = Dauer-Drift. Bekannter Schmerz im TF-Provider.
- **Kusto-Namen vs. Kubernetes-Namen.** Tabellen heißen `RawEvents`, K8s-Namen sind lowercase DNS-1123. Der ADX-Name muss aus der `crossplane.io/external-name`-Annotation oder einem Spec-Feld kommen, nicht aus `metadata.name` (siehe F11).
- **Async-Operationen.** `.create async materialized-view ... backfill=true` liefert eine Operation-ID; Backfill dauert Stunden. Observe muss `.show operations <id>` pollen und "Creating" sauber abbilden.
- **Fehlerklassifikation.** `.show table X` auf nicht existierende Tabelle = HTTP 400 `BadRequest_EntityNotFound` → muss zu `ResourceExists=false`, nicht zu einem Reconcile-Error. Alle anderen 400er sind echte Fehler. Braucht eine zentrale Fehler-Mapping-Schicht.
- **Vererbung von Policies.** Tabellen-Retention `null` heißt "erbt von Database". `.delete table X policy retention` setzt auf Vererbung zurück. Eine Policy-MR zu löschen = zurück auf Default, nicht "Entität löschen" – muss so dokumentiert sein.
- **Rechte des Provider-Principals.** Database Admin für alles auf DB-Ebene; Cluster `AllDatabasesAdmin` für Cluster-Policies, Workload Groups, Cluster-Principals. Das gehört in die Doku und bestimmt das ProviderConfig-Design (F5).

## 2. Artefakt-Inventar (Vorschlag für Priorisierung)

**Tier 1 (MVP, Datenbank-Ebene):**
Table, Function, MaterializedView, IngestionMapping (csv/json/avro/parquet/orc/w3clog), ExternalTable (Storage, Delta; SQL später), ContinuousExport, SecurityRole (database/table/function/materialized-view/external-table Principals).
Policies (Table / MaterializedView / Database, je nach Gültigkeit): Retention, Caching, Update, RowLevelSecurity, IngestionBatching, StreamingIngestion, Merge, Sharding, Partitioning, IngestionTime, AutoDelete, RestrictedViewAccess, ExtentTagsRetention, Encoding (column-level), ManagedIdentity (database-level).

**Tier 2 (Cluster-Ebene, braucht AllDatabasesAdmin):**
WorkloadGroup, RequestClassificationPolicy, CalloutPolicy, CapacityPolicy, SandboxPolicy, QueryWeakConsistencyPolicy, ManagedIdentityPolicy (cluster), MultiDatabaseAdmins, ClusterSecurityRole.

**Tier 3 (neuer / Nische):**
EntityGroup, GraphModel, QueryAccelerationPolicy (External Delta Tables), RowOrderPolicy, MirroringPolicy (Fabric), ExternalTable auf SQL/Cosmos/MySQL/Postgres.

**Explizit raus:** alles, was Daten bewegt (`.ingest`, `.set-or-append`, Stored Query Results), alles ARM (siehe P3), Ad-hoc-Queries.

## 3. Offene Fragen

Format: Frage → meine Empfehlung → Platz für deine Antwort. Wo ich eine klare Meinung habe, steht sie da; wo nicht, sage ich es.

### A. Scope

**F1 – ARM-Ressourcen raus?** Cluster, Database, Data Connections, Principal Assignments, Scripts, Attached DBs bleiben bei provider-upbound-azure; wir liefern nur Kusto-Endpoint-Artefakte.
Empfehlung: Ja, klar raus. Duplizieren wäre Pflege ohne Nutzen.
Antwort: **Ja, raus.** Nur Artefakte hinter dem Kusto-Endpoint. ARM bleibt bei provider-upbound-azure.

**F2 – Microsoft Fabric Eventhouse / KQL Database als Ziel?** Gleiche Engine, Endpoint `*.kusto.fabric.microsoft.com`, gleiche Entra-Tokens. Sollte funktionieren, ist aber ungetestet und Fabric hat Abweichungen (z. B. keine Cluster-Policies, Rechte-Modell über Fabric-Workspace).
Empfehlung: "Best effort, nicht getestet" in v1. Keine Fabric-spezifischen Features bauen.
Antwort: **Best effort, ungetestet.** Keine Fabric-spezifischen Features in v1, Doku weist darauf hin.

**F3 – Breite zuerst oder Tiefe zuerst?** Der Text oben sagt "komplett ausschmücken". Ich halte das für die falsche Reihenfolge: Erst ein vertikaler Schnitt (ProviderConfig + Auth + Table + Function + UpdatePolicy + RetentionPolicy) mit vollständigem Observe/Create/Update/Delete, Emulator-Tests und einem gebauten xpkg. Dann das Muster auf die restlichen Artefakte ausrollen. Sonst haben wir 40 halbfertige Controller und P1/P2 in jedem einzeln.
Empfehlung: Tiefe zuerst, dann Tier 1 komplett, dann Tier 2/3 nach Bedarf.
Antwort: **Tiefe zuerst.** M1 = Table, Function, UpdatePolicy, RetentionPolicy mit vollständigem Lifecycle, Tests und xpkg. Danach Muster ausrollen.

**F4 – Welches Artefakt brauchst du als erstes produktiv?** Das bestimmt die Reihenfolge innerhalb von Tier 1 und welcher Test-Cluster nötig ist.
Antwort: **Alle vier Gruppen** (Table/Function/UpdatePolicy, Policies, MaterializedView, SecurityRole/ExternalTable/ContinuousExport) werden produktiv gebraucht. Keine Priorisierung durch Bedarf → Reihenfolge innerhalb Tier 1 nach Risiko: erst die Kinds, die die schwierigen Muster erzwingen (Drift-Normalisierung, async, Principal-Normalisierung), damit die Basis früh steht.

### B. API-Design

**F5 – Wo lebt die Cluster-URI?** Option a) in der ProviderConfig (ein PC pro Cluster, wie provider-sql/provider-helm/provider-kubernetes). Option b) in `spec.forProvider.clusterUri` jeder MR (ein PC pro Identität). Bei a) muss eine Composition, die Cluster + Tabellen anlegt, auch eine PC erzeugen – das ist etabliert.
Empfehlung: a). Database-Name dagegen in den Spec jeder DB-scoped MR, nicht in die PC (sonst PC-Explosion).
Antwort: **ProviderConfig pro Cluster.** PC = Identität + Cluster-URI. Database-Name im Spec jeder DB-scoped MR.

**F6 – Policies als eigene Kinds oder inline in Table?** Und wenn eigene Kinds: ein Kind pro Policy mit Ziel-Diskriminator (`entity: {kind: Table|MaterializedView|Database, name}`) oder pro Kombination wie im TF-Provider (`TableRetentionPolicy`, `MaterializedViewRetentionPolicy`, …)? Letzteres wären ~45 CRDs nur für Policies.
Empfehlung: Eigene Kinds (entspricht dem Kusto-Modell: unabhängig alter-/löschbar, Vererbung), ein Kind pro Policy-Typ mit Ziel-Diskriminator und CEL-Validierung, welche Ziele erlaubt sind.
Antwort: **Ein Kind pro Policy-Typ** mit Ziel-Diskriminator `entity: {kind: Table|MaterializedView|Database, name}`. CEL validiert erlaubte Ziele pro Policy.

**F7 – Table-Schema als strukturierte Liste oder KQL-String?** `columns: [{name, type, docstring}]` ist validierbar (Typ-Enum per CEL) und diffbar. Der TF-Provider bietet zusätzlich `table_schema: "a:string,b:int"` und `from_query`.
Empfehlung: Nur strukturierte Liste. `from_query` ist Ingestion, nicht Schema-Deklaration → raus (YAGNI).
Antwort: **Strukturierte Spaltenliste** `columns: [{name, type, docstring}]`. Kein KQL-String, kein from_query.

**F8 – Verhalten bei destruktiven Änderungen (P2).** Vorschlag: `schemaUpdateMode: Merge|Replace` (Default `Merge` = `.alter-merge`, nur Hinzufügen). Spaltentyp-Änderung → immer `Synced=False` mit klarer Meldung, nie automatisch. MV-Query-Änderung → `.alter materialized-view` versuchen, bei Fehler Status-Meldung, kein Auto-Recreate. Löschen → Crossplane-Standard (`deletionPolicy`/`managementPolicies`), keine Zusatz-Guards. Frage: Willst du für Table einen abweichenden Default (Orphan statt Delete) als Datenverlust-Schutz? Das widerspricht der Crossplane-Konvention und überrascht Nutzer.
Empfehlung: Crossplane-Standard beibehalten, Doku warnt deutlich; `schemaUpdateMode` wie beschrieben.
Antwort: **Merge-Default, Standard-Delete.** `schemaUpdateMode: Merge|Replace`, Default Merge. Typänderung → Synced=False mit Meldung, nie automatisch. Löschen nach Crossplane-Standard, Doku warnt deutlich.

**F9 – Crossplane v2 only (namespaced MRs) oder dual (zusätzlich cluster-scoped legacy)?** Dual verdoppelt CRDs und Testmatrix. Cluster-scoped wird abgekündigt.
Empfehlung: v2-only, namespaced, Mindestversion Crossplane 2.0.
Antwort: **v2 only, namespaced MRs.** Mindestversion Crossplane 2.0. Keine cluster-scoped Legacy-Kinds.

**F10 – Cross-Resource-Referenzen (`tableRef`/`tableSelector`)?** Zwischen eigenen Kinds möglich (UpdatePolicy → Table, ContinuousExport → ExternalTable, MV → Source Table). Referenz auf eine upbound-azure `Database`-MR ist technisch nicht möglich (fremdes Scheme) → Database bleibt String, wird per Composition gepatcht.
Empfehlung: Refs für Table/Function/ExternalTable/MaterializedView in v1, generiert per crossplane-tools. Kostet Boilerplate, spart Reconcile-Fehler durch falsche Reihenfolge.
Antwort: **Ja, in v1.** Refs/Selectors für Table, Function, ExternalTable, MaterializedView, generiert per crossplane-tools. Database bleibt String.

**F11 – Woher kommt der ADX-Name?** Crossplane-Konvention: `crossplane.io/external-name`-Annotation, Default `metadata.name`. Problem: `metadata.name` ist lowercase, ADX-Namen sind meist PascalCase → Nutzer müssten immer die Annotation setzen.
Empfehlung: Konvention einhalten (Annotation ist Source of Truth für Import), zusätzlich optionales `spec.forProvider.name`, das beim ersten Reconcile in die Annotation geschrieben wird. Vorher prüfen, ob das gegen aktuelle Crossplane-Guidelines verstößt – da bin ich nicht sicher.
Antwort: **Annotation + optionales `spec.forProvider.name`.** external-name-Annotation bleibt Source of Truth; Spec-Feld wird beim ersten Reconcile in die Annotation übernommen. Vorher gegen aktuelle Crossplane-Guidelines prüfen.

**F12 – Security Roles: autoritativ oder additiv?** Eine MR = komplette Principal-Liste für (Entität, Rolle) via `.set` (autoritativ, entfernt Fremde). Oder eine MR pro Principal via `.add`/`.drop` (additiv, mehrere Teams können dieselbe Rolle bestücken, aber kein Schutz gegen manuell hinzugefügte Principals).
Empfehlung: Beides als `mode: Authoritative|Additive` am Kind, Default `Authoritative`. Plus Principal-Normalisierung (siehe Stolpersteine) von Anfang an.
Antwort: **Beides, Default Authoritative.** `mode: Authoritative|Additive` am Kind. Principal-Normalisierung (Object-ID vs. App-ID, PrincipalFQN) von Anfang an.

**F13 – External Tables mit Secrets.** Connection String per Secret-Ref, wird immer als `h@'...'` gesendet → kein Drift-Vergleich auf dem Connection String möglich. Akzeptabel? Alternative wäre, Klartext zu senden, damit `.show` vergleichbar ist – dann steht das Secret für jeden Database-Viewer sichtbar in der Metadata.
Empfehlung: Verschleiern, Write-only, Doku sagt es klar. Managed Identity als empfohlener Pfad.
Antwort: **Verschleiern, write-only.** Connection String per Secret-Ref, immer als `h@'...'` gesendet. Kein Drift-Vergleich auf dem String, Doku sagt es klar. Managed Identity als empfohlener Pfad.

**F14 – Materialized View Backfill in v1?** Braucht async Create + Operation-Polling + "Creating"-Zustand über Stunden.
Empfehlung: Ja, aber async von Anfang an (sync würde Timeouts produzieren). Ohne Backfill wäre der Kind für Bestandsdaten nutzlos.
Antwort: **Ja, async von Anfang an.** `.create async materialized-view`, Operation-ID im Status, Polling über `.show operations`, Zustand Creating sauber abgebildet.

**F15 – Was gehört in `status.atProvider`?** Minimal (observierter Zustand für Debugging) oder auch Betriebsdaten (MV: `IsHealthy`, `LastRun`, `MaterializedTo`; ContinuousExport: `LastRunResult`, `ExportedTo`)?
Empfehlung: v1 minimal; Betriebsdaten sind billig hinzuzufügen, wenn jemand sie braucht.
Antwort: **Minimal.** Observierter Zustand für Debugging. Betriebsdaten (MV-Health, Export-Status) später additiv.

### C. Authentifizierung & Betrieb

**F16 – Welche Auth-Methoden in v1?** Client Secret (K8s Secret), Client Certificate, Workload Identity (AKS federated), Managed Identity (system/user assigned), Azure CLI (nur Dev).
Empfehlung: Client Secret + Workload Identity + Managed Identity in v1 (alles via `azidentity`, wenig Mehraufwand). Certificate später.
Antwort: **Client Secret, Workload Identity, Managed Identity** (alle via azidentity). Client Certificate nicht in v1.

**F17 – Sovereign Clouds (China, US Gov)?** Andere Endpoints und Token-Authority.
Empfehlung: Ein Feld `azureEnvironment` in der PC (billig, `azidentity` kann es). Nicht testen.
Antwort: **Feld `azureEnvironment` in der ProviderConfig, ungetestet.** Public Cloud als Default.

**F18 – Größenordnung?** Wie viele MRs pro Cluster erwartest du: 50, 500, 5.000? Ab ~500 wird P4 zum Designthema (Batch-Observe per `.show database schema as json`, Caching), unter 100 reicht ein sinnvolles Poll-Intervall.
Antwort: **500 bis 5.000 MRs pro Cluster.** Konsequenz: Batch-Observe (z. B. `.show database schema as json` / `.show database policies` pro Database mit kurzem Cache) wird Designthema in M0, nicht Nachrüstung. Bei 5.000 MRs und 10-Minuten-Poll sind das ~8 Commands/s Dauerlast ohne Batching. Throttling-Handling (429/Backoff) ist Pflicht.

**F19 – Poll-Intervall-Default?** Crossplane-Default 1 Minute, Upbound-Provider 10 Minuten.
Empfehlung: 10 Minuten, per Flag änderbar.
Antwort: **10 Minuten**, per Flag änderbar.

### D. Projekt, Tests, Veröffentlichung

**F20 – Name, API-Group, Org, Lizenz.** `provider-azure-adx` oder `provider-kusto`? API-Group braucht eine Domain, die dir gehört oder crossplane-contrib (z. B. `adx.azure.crossplane.io` nur, wenn das Projekt nach crossplane-contrib soll). Lizenz Apache-2.0 wie Crossplane? Ziel crossplane-contrib oder eigene Org?
Empfehlung: Apache-2.0. Name `provider-azure-adx` (Nutzer suchen nach "ADX"), API-Group unter eigener Domain, Wechsel zu crossplane-contrib später ist ein Breaking Change der Group → vorher entscheiden.
Antwort: **Name `provider-azure-adx`, Lizenz Apache-2.0, Ziel crossplane-contrib.** API-Group damit unter `crossplane.io` (Vorschlag `adx.azure.crossplane.io`; prüfen, ob contrib-Provider mit namespaced Kinds das `.m.`-Infix erwarten wie provider-helm/provider-kubernetes). Risiko: Wird die Aufnahme abgelehnt, muss die Group vor dem ersten öffentlichen Release umbenannt werden → Aufnahme-Antrag früh stellen (Prozess über crossplane/org prüfen), bis dahin kein v0.1.0 öffentlich. Contrib verlangt DCO, CODE_OF_CONDUCT, OWNERS mit mindestens einem committed Maintainer. **Hosting (2026-09-07):** zunächst `github.com/FakieHeelflip/provider-azure-adx`, Images unter `ghcr.io/fakieheelflip/…`, Transfer nach crossplane-contrib nach Aufnahme. **Revidiert 2026-09-07 (abends):** Der Provider bleibt bei functional.team: Repo `github.com/functional-team/provider-azure-adx`, Images `ghcr.io/functional-team/provider-azure-adx`, API-Groups unter `functional.team` (`adx.functional.team`, `policy.adx.functional.team`, `security.adx.functional.team`, `cluster.adx.functional.team`). Kein Transfer nach crossplane-contrib geplant; damit entfällt die `.m.`-Frage (S9). Eine spätere Aufnahme in contrib wäre ein Group-Wechsel und damit ein Breaking Change.

**F21 – Test-Strategie und Budget.** Unit-Tests mit gemocktem Kusto-Client; Integrationstests gegen den Emulator (Testcontainers) für alles, was er kann; E2E gegen einen echten Cluster für Security Roles, External Tables auf Storage, Continuous Export, Managed Identity. Letzteres braucht eine Azure-Subscription, einen SPN mit Cluster-Admin-Rechten und laufende Kosten (kleinster Dev-Cluster ~ 100 €/Monat, Start/Stop per Pipeline möglich).
Frage: Gibt es eine Subscription/Budget dafür? Wenn nein, welche Artefakte akzeptieren wir als "ungetestet gegen echten Cluster"?
Antwort: **Ja, Subscription und Budget vorhanden.** E2E gegen echten Dev-Cluster ab M2, Start/Stop per Pipeline. SPN mit Cluster-AllDatabasesAdmin und Storage-Account für External-Table-/Continuous-Export-Tests einplanen.

**F22 – Registry & Release.** xpkg nach `xpkg.upbound.io` (Marketplace-Listing) oder `ghcr.io`? Release-Tooling: Crossplane `build`-Submodule (Standard) mit GitHub Actions; Renovate für Dependencies.
Empfehlung: ghcr.io + build-Submodule, Marketplace optional später.
Antwort: **ghcr.io + Crossplane build-Submodule**, GitHub Actions, Renovate. Nach contrib-Aufnahme zusätzlich xpkg.crossplane.io/crossplane-contrib.

**F23 – Doku-Umfang.** Pro Kind ein Beispiel-YAML (Pflicht für `make reviewable`), generierte CRD-Referenz, ein Kapitel "Migration vom TF-Provider" mit Feld-Mapping.
Empfehlung: Ja zu allen drei, Migrationskapitel aber erst, wenn Tier 1 stabil ist.
Antwort: **Alle drei.** Beispiel-YAML pro Kind, generierte CRD-Referenz, Migrationskapitel vom TF-Provider nach stabilem Tier 1.

**F24 – Ist der alte Ordner `provider-azure-adx-old` löschbar?** Enthält keinen eigenen Code (siehe oben).
Empfehlung: Löschen.
Antwort: **Gelöscht** am 2026-09-07.

## 4. Entscheidungslog (2026-09-07)

Alle 24 Fragen sind beantwortet, Details stehen bei den Fragen oben. Kurzfassung:

- **Scope:** nur Artefakte hinter dem Kusto-Endpoint, kein ARM. Fabric Eventhouse best effort, ungetestet.
- **Vorgehen:** Tiefe zuerst. Alle vier Tier-1-Gruppen werden produktiv gebraucht, daher Reihenfolge nach Risiko statt nach Bedarf.
- **API:** ProviderConfig pro Cluster, Database im Spec. Ein Kind pro Policy-Typ mit `entity`-Diskriminator. Strukturierte Spaltenliste. `schemaUpdateMode: Merge|Replace`, Default Merge, Typänderung nie automatisch. Crossplane-Standard-Delete. v2-only, namespaced. Refs/Selectors in v1. external-name-Annotation plus optionales `spec.forProvider.name`. SecurityRole `mode: Authoritative|Additive`. External-Table-Secrets write-only und verschleiert. MaterializedView async mit Backfill. `status.atProvider` minimal.
- **Auth & Betrieb:** Client Secret, Workload Identity, Managed Identity. Feld `azureEnvironment`. 500 bis 5.000 MRs pro Cluster erwartet, daher Batch-Observe ab M0. Poll-Default 10 Minuten.
- **Projekt:** `provider-azure-adx`, Apache-2.0, Ziel crossplane-contrib mit Group unter `crossplane.io` (`.m.`-Infix-Frage klären, Antrag früh stellen). ghcr.io + build-Submodule. Doku: Beispiele, CRD-Referenz, Migrationskapitel später. Echter Test-Cluster mit Budget vorhanden. Alter Ordner gelöscht.

## 5. Meilensteine (aktualisiert nach den Antworten)

- **M0 – Fundament:** provider-template v2 (namespaced). ProviderConfig mit Client Secret / Workload Identity / Managed Identity und `azureEnvironment`. Kusto-Client-Wrapper auf `azkustodata` mit zentralem Fehler-Mapping (EntityNotFound → nicht vorhanden), 429-Backoff, Quoting-Helfern. **Batch-Observe-Design** wegen F18: Observe pro Database gebündelt (`.show database schema as json`, `.show database ... policies`) mit kurzem Cache, statt ein `.show` pro MR. Emulator per Testcontainers in CI. Contrib-Vorbereitung: DCO, CODE_OF_CONDUCT, OWNERS, Group-Frage klären, Antrag stellen.
- **M1 – Vertikaler Schnitt:** Table (`schemaUpdateMode`, Name über Annotation/Spec), Function (Body- und Parameter-Normalisierung), UpdatePolicy, RetentionPolicy (`entity`-Diskriminator). Refs Table ↔ UpdatePolicy/Function. Vollständige Lifecycle-Tests im Emulator. Erstes internes xpkg, noch nicht öffentlich, solange die Group nicht entschieden ist.
- **M2 – Tier 1 komplett:** restliche Policies, MaterializedView (async, Backfill, Operation-Polling), IngestionMapping, SecurityRole (Authoritative|Additive, Principal-Normalisierung), ExternalTable (write-only Secrets, Managed-Identity-Pfad), ContinuousExport. E2E gegen echten Dev-Cluster mit Start/Stop-Pipeline, SPN mit AllDatabasesAdmin, Storage-Account für Export-Tests. Erstes öffentliches Release v0.1.0.
- **M3 – Tier 2:** Cluster-Ebene (WorkloadGroup, RequestClassification, Callout, Capacity, Sandbox, ManagedIdentity, ClusterSecurityRole).
- **M4 – Tier 3 nach Bedarf.** Migrationskapitel TF → Crossplane nach M2.

Abbruchkriterium für M1: Wenn P1 (Drift-Normalisierung) für Functions nicht stabil zu lösen ist, muss das Design auf "Hash-Vergleich in einer Annotation statt Live-Vergleich" umgestellt werden. Das ist ein Kompromiss (kein Drift-Schutz gegen manuelle Änderungen im Cluster) und sollte bewusst entschieden werden.
