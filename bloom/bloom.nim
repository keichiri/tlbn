import md5, os, bignum



proc hashF(input: string): string =
    $toMD5(input)


proc hashToBigInt(hash: string): bignum.Int =
    bignum.newInt(hash, 16)



type
    BloomFilter = ref object
        bitfield: seq[byte]
        hashFunc: (proc(item: string): string)
        rounds: int


proc newBloomFilter(length: int, rounds: int, hashFunc: (proc(item: string): string)): BloomFilter =
    let byteArrayLength: int = length div 8
    let byteArray = newSeq[byte](byteArrayLength)
    result = BloomFilter(
        bitfield: byteArray,
        hashFunc: hashFunc,
        rounds: rounds,
    )


proc add(bf: BloomFilter, item: string) =
    var input = item
    for _ in 1 .. bf.rounds:
        input = bf.hashFunc(input)
        let hashValue = hashToBigInt(input)
        let bitIndex = (hashValue mod (bf.bitfield.len() * 8)).toInt()
        let byteIndex = bitIndex div 8
        let byteOffset = bitIndex mod 8
        let value = byte(1 shl (7 - byteOffset))
        bf.bitfield[byteIndex] = bf.bitfield[byteIndex] or value


proc isMember(bf: BloomFilter, item: string): bool =
    var input = item
    for _ in 1 .. bf.rounds:
        input = bf.hashFunc(input)
        let hashValue = hashToBigInt(input)
        let bitIndex = (hashValue mod (bf.bitfield.len() * 8)).toInt()
        let byteIndex = bitIndex div 8
        let byteOffset = bitIndex mod 8
        let value = byte(1 shl (7 - byteOffset))
        if (bf.bitfield[byteIndex] and value) == 0:
            return false

    return true




proc getInput(filePath: string): seq[string] =
    let file = open(filePath, fmRead)
    result = @[]
    for line in file.lines:
        result.add(line)


when isMainModule:
    proc demo() =
        let words = getInput("/usr/share/dict/words")
        let bloomFilter = newBloomFilter(4000000, 7, hashF)
        var toadd = newSeq[string](0)
        var toskip = newSeq[string](0)
        for i, word in words:
            if i mod 100 == 0:
                toskip.add(word)
            else:
                toadd.add(word)

        for word in toadd:
            bloomFilter.add(word)

        var falseNegatives = 0
        for word in toadd:
            if not bloomFilter.isMember(word):
                falseNegatives += 1

        echo("False negatives: ", falseNegatives)


        var falsePositives = 0
        for word in toskip:
            if bloomFilter.isMember(word):
                falsePositives += 1

        echo("To skip: ", toskip.len())
        echo("False positives: ", falsePositives)

    demo()


