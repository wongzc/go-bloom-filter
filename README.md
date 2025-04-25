- Using go to implement the concept of bloom filter
- Bloom filter is, in simple word, array+ multiple hash functions to check if something presented or not
- Hash functions: to compute hash values of a target
- Array: to stored the result of hash functions
- When checking, we compute the hash values using hash function
    - then if all bit of hash values present in array: maybe present
    - if any bit not present: 100% not present
- So bloom filter can use for:
    - case when likely return false
    - check malicious URL
    - check if user have read a article
- Advantage
    - fast and memory efficient ( just need 1 array to save so many thing!)

- Formula 

$$
m = -\frac{n \cdot \ln p}{(\ln 2)^2}
$$

$$
k = \frac{m}{n} \cdot \ln 2
$$

    - m: bit array size  
    - n: number of item that going to save  
    - p: acceptable false positive rate  
    - k: number of hash function  

- Feature:
    - use xxHash for optimum speed, collision rate, distribution
    - use byte slices instead of bool slice to compress memory usage to 1/8
    - use double hashing of k times with index to achieve k of hash values
    - dynamiccaly set slice size and hash function k by user required acceptable false value rate and exepected item store