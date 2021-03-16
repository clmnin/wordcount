import sys

counts = {}
for line in sys.stdin:
    words = line.lower().split()
    for word in words:
        counts[word] = counts.get(word, 0) + 1

pairs = sorted(counts.items(), key=lambda kv: kv[1], reverse=True)
for word, count in pairs:
    print(word, count)
