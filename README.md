- using go to implement the concept of bloom filter
- bloom filter is, in simple word, a array+ multiple hash functions
- the hash functions to compute hash value of a target
- the array used to stored the result of hash functions
- when checking, we compute the hash values using hash function
    - then if all bit of hash values present in array: maybe present
    - if any bit not present: 100% not present
- so bloom filter can use for:
    - case when likely return false
    - check malicious URL
    - check if user have read a article
- advantage
    - fast and memory efficient ( u just need 1 array to save so many thing!)

- formula 
    - m: bit array size  
    - n: number of item that going to save  
    - p: acceptable false positive rate  
    - k: number of hash function  
$$
m = -\frac{n \cdot \ln p}{(\ln 2)^2}
$$

$$
k = \frac{m}{n} \cdot \ln 2
$$
