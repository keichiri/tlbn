import sys
import time
from hashlib import sha256, md5, sha1
from math import log, ceil
from random import random


class BloomFilter:
    @staticmethod
    def calculate_parameters(input_size, false_positive_ratio):
        bit_length = -input_size * log(false_positive_ratio) / (log(2) ** 2)
        rounds = bit_length / input_size * log(2)
        return ceil(bit_length), ceil(rounds)

    bit_weights = [128, 64, 32, 16, 8, 4, 2, 1]

    def __init__(self, length, hash_function, rounds):
        self._bitfield = bytearray(length // 8)
        # Cannot use the 'length' parameter since it might be longer than len(self._bitfield) * 8 after division
        self._length = len(self._bitfield) * 8
        self._hash_function = hash_function
        self._rounds = rounds

    def add(self, element):
        hash_input = element
        for _ in range(self._rounds):
            hash_input = self._hash_function(hash_input)
            bit_pos = int.from_bytes(hash_input, 'big', signed=False) % self._length
            byte_pos, byte_offset = bit_pos // 8, bit_pos % 8
            self._bitfield[byte_pos] |= self.bit_weights[byte_offset]

    def is_member(self, element):
        hash_input = element
        for _ in range(self._rounds):
            hash_input = self._hash_function(hash_input)
            bit_pos = int.from_bytes(hash_input, 'big', signed=False) % self._length
            byte_pos, byte_offset = bit_pos // 8, bit_pos % 8
            if not self._bitfield[byte_pos] & self.bit_weights[byte_offset]:
                return False

        return True


def hash_func(input):
    hasher = sha1()
    hasher.update(input)
    return hasher.digest()


def demo():
    ratio = 0.01
    with open('/usr/share/dict/words') as f:
        items = f.read().split()

    start = time.time()
    items = [item.encode() for item in items]
    to_add = []
    to_check = []

    for item in items:
        if random() < ratio:
            to_check.append(item)
        else:
            to_add.append(item)

    optimal_len, optimal_rounds = BloomFilter.calculate_parameters(len(to_add), ratio)
    print(f'Item count: {len(to_add)}. Optimal bit length: {optimal_len}. Optimal rounds: {optimal_rounds}')
    bf = BloomFilter(optimal_len, hash_func, optimal_rounds)
    for item in to_add:
        bf.add(item)

    for item in to_add:
        if not bf.is_member(item):
            raise Exception

    false_positives = 0
    for item in to_check:
        if bf.is_member(item):
            false_positives += 1
    end = time.time()

    print(f'False positives: {false_positives}. False positive ratio: {float(false_positives) / len(to_check)}. Elapsed time: {end-start}')

    s = {item for item in to_add}
    d = {item: True for item in to_add}
    print(f'Total memory used by filter bitfield: {sys.getsizeof(bf._bitfield)}')
    print(f'Total memory used by set: {sys.getsizeof(s)}')
    print(f'Total memory used by dictionary: {sys.getsizeof(d)}')
    print(f'Total memory used by list: {sys.getsizeof(to_add)}')


if __name__ == '__main__':
    demo()
