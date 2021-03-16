# Count Words Problem

A good interview question - https://benhoyt.com/writings/count-words/

## Problem statement and constraints

From the blog post

```markdown
Each program must read from standard input and print the frequencies of unique, space-separated words, in order from most frequent to least frequent. To keep our solutions simple and consistent, here are the (self-imposed) constraints I’m working against:

* Case: the program must normalize words to lowercase, so “The the THE” should appear as “the 3” in the output.
* Words: anything separated by whitespace – ignore punctuation. This does make the program less useful, but I don’t want this to become a tokenization battle.
* ASCII: it’s okay to only support ASCII for the whitespace handling and lowercase operation. Most of the optimized variants do this.
* Ordering: if the frequency of two words is the same, their order in the output doesn’t matter. I use a normalization script to ensure the output is correct.
* Threading: it should run in a single thread on a single machine (though I often discuss concurrency in my interviews).
* Memory: don’t read whole file into memory. Buffering it line-by-line is okay, or in chunks with a maximum buffer size of 64KB. That said, it’s okay to keep * the whole word-count map in memory (we’re assuming the input is text in a real language, not full of randomized unique words).
* Text: assume that the input file is text, with “reasonable” length lines shorter than the buffer size.
* Safe: even for the optimized variants, try not to use unsafe language features, and don’t drop down to assembly.
* Hashing: don’t roll our own hash table (with the exception of the optimized C version).
* Stdlib: only use the language’s standard library functions.

Our test input file will be the text of the King James Bible, concatenated ten times. I sourced this from Gutenberg.org, replaced smart quotes with the ASCII quote character, and used cat to multiply it by ten to get the 43MB reference input file.
```

## Implementations

We'll implement the solution in different languages and compare their results

| Language | Simple |Optimized | Notes |
|---|---|---|---|
| grep | 0.04 | 0.04 | grep reference; optimized sets LC_ALL=C
| wc -w| 0.29 | 0.20 | wc reference; optimized sets LC_ALL=C
| C | 0.97 | 0.23 | |
| Go | 1.14 | 0.38 | |
| Rust A | 1.41 | 0.35 | by Andrew Gallant |
| Rust B | 1.48 | 0.28 | also by Andrew: bonus and custom hash |
| C++ | 1.75 | 0.98 | “optimized” isn’t very optimized |
| Python | 2.07 | 1.27 | |
| C#	| 3.43	|  	| original by John Taylor|
| AWK	| 3.52	| 1.13	|optimized uses mawk |
| Forth	| 4.21	| 1.44 | |
| Shell	| 14.67	| 1.86 | optimized does LC_ALL=C sort -S 2G|

### Python 

1. Most basic implementation

    This is a beginner approach to the problem. It's a beginner code becase they fail to use the python tools a pythonista would use
    ```sh
    cat input.txt | python basic.py
    ```

2. Simple

    An idiomatic Python version would probably use collections.Counter. Python’s collections library is really nice

    This is Unicode aware and probabily what most will write in real life. It’s actually quite efficient, because all the low-level stuff is really done in C

    Pros:
    * Low level stuff done in C
        * reading file
        * converting to lowercase
        * splitting on whitespace
        * updating counter
        * the sorting that Counter.most_common does

    ```sh
    python simple.py < input.txt
    ```

    **Profile the code**
    ```
    python -m cProfile -s tottime simple.py < input.txt
    ```
    
    **Observations**
    
    * 998,170 is the number of lines in the input, and because we’re reading line-by-line, we’re calling functions and executing the Python loop that many times.
    * The large amount of time spent in simple.py itself shows how (relatively) slow it is to execute Python bytecode – the main loop is pure Python, again executed 998,170 times.
    * str.split is relatively slow, presumably because it has to allocate and copy many strings.
    * Counter.update calls isinstance, which adds up. I thought about calling the C function _count_elements directly, but that’s an implementation detail and I decided it fell into the “unsafe” category.

3. Optimized
    Let's read it in 64KB chunks.

    Pros:
    * Instead of our main loop processing 42 characters at a time (the average line length), we’re processing 65,536 at a time (less the partial line at the end).
    * We’re still reading and processing the same number of bytes, but we’re now doing most of it in C rather than in the Python loop.
    * process things in bigger chunks - best way to optimize in python. Let C do it's work


    ```sh
    python optimized.py < input.txt
    ```

    **Profile the code**
    ```
    python -m cProfile -s tottime optimized.py < input.txt
    ```

    **Observations**

    * The _count_elements and str.split functions are still taking most of the time, but they’re only being called 662 times instead 998170 (on roughly 64KB at a time rather than 42 bytes)

### Go

1. Simple
    A simple, idiomatic Go version would probably use bufio.Scanner with ScanWords as the split function. 

    ```sh
    go build -o simple-go simple.go
    ./simple-go < input.txt
    ```

    **Profile the code**
    We have added the profiled inside the code

    run this to read the profile of the code
    ```sh
    go tool pprof -http=:7777 cpuprofile_simple
    ```

    **Observations**

    * the operations in the per-word hot loop take all the time.
    * A good chunk of the time is spent in the scanner, and another chunk is spent allocating strings to insert into the map

2. Optimizied
    * To improve scanning, we’ll essentially make a cut-down version of bufio.Scanner and ScanWords (and do an ACIII to-lower operation in place).
    * To reduce the allocations, we’ll use a `map[string]*int` instead of `map[string]int` so we only have to allocate once per unique word, instead of for every increment

    ```sh
    go build -o optimized-go optimized.go
    ./optimized-go < input.txt
    ```

    **Profile the code**
    We have added the profiled inside the code

    run this to read the profile of the code
    ```sh
    go tool pprof -http=:7777 cpuprofile_optimized
    ```

    **Observations**

    * Go gives you a fair bit of low-level control (and you could go quite a lot further – memory mapped I/O, a custom hash table, etc)
    * A lot of thought and time to write this though
    * In practice one would stick to bufio.Scanner with ScanWords, bytes.ToLower, and the map[string]*int trick.

### C++

1. Simple

    ```sh
    g++ -O2 simple.cpp -o simple-cpp
    ./simple-cpp < input.txt
    ```

    **Profile the code**
    I don't know

    **Observations**
    * the first thing to do is compile with optimizations enabled (g++ -O2).
    * noticed that I/O was comparatively slow
    * It turns out there is a magic incantation you can recite at the start of your program to disable synchronizing with the C stdio functions after each I/O operation. This line makes it run almost twice as fast: `ios::sync_with_stdio(false);`

2. Optimized
    ```sh
    g++ -O2 -DNDEBUG  -std=c++17 optimized.cpp -o optimized-cpp
    ./optimized-cpp < input.txt
    ```

    **Profile the code**
    I don't know
    
    **Observations**
    * Further optimizations is writing more low-level code. Which will be more C code

### C
C is a beautiful beast that will never die: fast, unsafe, and simple. It’s also ubiquitous (the Linux kernel, Redis, PostgreSQL, SQLite, many many libraries … the list is endless), and it’s not going away anytime soon.

Unfortunately, C doesn’t have a hash table data structure in its standard library. However, there is libc, which has the hcreate and hsearch hash table functions, so we’ll make a small exception and use those libc-but-not-stdlib functions. In the optimized version we’ll roll our own hash table.

One minor annoyance with hcreate is you have to specify the maximum table size up-front. I know the number of unique words is about 30,000, so we’ll make it 60,000 for now.

1. Simple
    ```sh
    gcc -O2 simple.c -o simple-c
    ./simple-c < input.txt
    ```

    **Profile the code**
    Ben Hoyt used Valgrind
    ```sh
    valgrind --tool=callgrind ./simple-c < input.txt >/dev/null
    ```
    
    **Observations**
    * There’s a fair bit of boilerplate (mostly for memory allocation and error checking), but as far as C goes, I don’t think it’s too bad.
    * The tricky stuff
        * tokenization behind scanf
        * hash table operations behind hsearch
    * relatively fast out of the box
    * performance
        * it shows that scanf is the major culprit
        * followed by hsearch
    
    **Improvements to be made**
    * Read the file in chunks, like we did in Go and Python. This will avoid the overhead of scanf.
    * Process the bytes only once, or at least as few times as possible – I’ll be converting to lowercase and calculating the hash as we’re tokenizing into words.
    * Implement our own hash table using the fast FNV-1 hash function.

2. Optimized
    ```sh
    gcc -O2 optimized.c -o optimized-c
    ./optimized-c < input.txt
    ```
    **Observations**
    * rolling your own hash table with linear probing is not a lot of code.
    * we haven't implemented a dynamic table size, but that's the best approach
    * it is small (17KB)
    * very fast
    * a little bit faster than the Go version, as we’ve rolled our own custom hash table, and we’re processing fewer bytes



    **Profile the code**
    Ben Hoyt used Valgrind
    ```sh
    valgrind --tool=callgrind ./simple-c < input.txt >/dev/null
    ```

### Rust
The Rust version was writen by Andrew Gallant or more commonly known as [BurntSushi](https://github.com/BurntSushi/ripgrep). He wrote a simple and an optimized version which had similar preformance to that of the Go implementation. But then he write 3 more variants. Let's check them out. 
1. Simple
    * is similar to the simple Go and C++ versions
2. Optimized
    * This version is an approximate port of the optimized Go program. 
    * Its buffer handling is slightly simpler: 
        * we don’t bother with dealing with the last newline character.
        * This may appear to save work, but it only saves work once per 64KB buffer, so is likely negligible. It’s just simpler IMO.
    * There’s nothing particularly interesting here other than swapping out std’s default hashing algorithm for one that isn’t cryptographically secure.
    * std uses a cryptographically secure hashing algorithm by default, which is a bit slower.
