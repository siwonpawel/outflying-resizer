# outflying-resizer

Multithreaded implementation of jpeg scaling aplication written in Go.

To run this tool:

> outflying-resizer -f /path/to/folder/with/images -o /path/to/output/folder -t thread_number -s scaling_in_percentage

Average time of 10 executions for different thread numbers. Boost in comparison to one thread

| Threads | Avergate time (s) | Boost (%) |
| :-----: | :---------------: | :-------: |
|    1    |      19,6959      |    ---    |
|    2    |      10,7795      |  45,270   |
|    3    |      8,0698       |  59,028   |
|    4    |      6,5027       |  66,984   |
|    5    |      5,5362       |  71,892   |
|    6    |      5,2229       |  73,482   |
|    7    |       5,059       |  74,314   |
|    8    |      4,9793       |  74,719   |
|    9    |      4,9549       |  74,843   |
|   10    |      4,5887       |  76,702   |
|   11    |      4,5496       |  76,901   |
|   12    |      4,5399       |  76,950   |