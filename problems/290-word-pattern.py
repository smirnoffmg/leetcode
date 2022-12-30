"""
Given a pattern and a string s, find if s follows the same pattern.

Here follow means a full match, such that there is a bijection between a letter in pattern and a non-empty word in s.
"""


class Solution:
    def wordPattern(self, pattern: str, s: str) -> bool:
        mapping = {}
        splitted = s.split()

        if len(pattern) != len(splitted):
            return False

        for i in range(len(pattern)):
            if pattern[i] not in mapping:
                if splitted[i] in mapping.values():
                    return False
                mapping[pattern[i]] = splitted[i]
            elif mapping[pattern[i]] != splitted[i]:
                return False

        return True


assert Solution().wordPattern("abba", "dog cat cat dog") == True
assert Solution().wordPattern("abba", "dog cat cat fish") == False
assert Solution().wordPattern("aaaa", "dog cat cat dog") == False
