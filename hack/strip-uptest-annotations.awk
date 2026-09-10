# Removes uptest hook annotations from an example manifest, together with the
# comment block that explains them and an annotations: key that is left with no
# children. Used when the examples are staged for the package: uptest hooks are
# test scaffolding, and the examples embedded in the package are what the
# Upbound Marketplace shows as the manifest for each kind.
#
# This is text surgery on purpose. uptest ships cleanupexamples, which does the
# same job by round-tripping every file through a YAML parser -- and drops all
# comments and the key order with it. The comments are the better half of these
# examples.
{ line[NR] = $0 }
END {
    for (i = 1; i <= NR; i++) {
        if (line[i] !~ /uptest\.upbound\.io\//) continue
        drop[i] = 1
        for (j = i - 1; j >= 1 && line[j] ~ /^[ \t]*#/; j--) drop[j] = 1
    }

    for (i = 1; i <= NR; i++) {
        if (drop[i] || line[i] !~ /^[ \t]*annotations:[ \t]*$/) continue
        match(line[i], /^[ \t]*/)
        indent = RLENGTH
        empty = 1
        for (j = i + 1; j <= NR; j++) {
            if (drop[j] || line[j] ~ /^[ \t]*$/) continue
            match(line[j], /^[ \t]*/)
            if (RLENGTH > indent) empty = 0
            break
        }
        if (empty) drop[i] = 1
    }

    for (i = 1; i <= NR; i++)
        if (!drop[i]) print line[i]
}
