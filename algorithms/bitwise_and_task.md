Bitwise AND of Numbers Range

Interesting LeetCode problem. The description of the problem is the following:

```
Given a range [m, n] where 0 <= m <= n <= 2147483647, return the bitwise AND of all numbers in this range, inclusive.

Example 1:

Input: [5,7]
Output: 4

Example 2:

Input: [0,1]
Output: 0
```

So, there are at least 2 solution of this probles (I'll write it down in my favorite programming language - Rust). First - simple, we will go from first number in range over the all presenten numbers (last number is included)

9   -- 0 0 0 0 1 0 0 1
10  -- 0 0 0 0 1 0 1 0
11  -- 0 0 0 0 1 0 1 1
12  -- 0 0 0 0 1 1 0 0

Итак, у нас числа. За O(n) решение простое. Просто присваиваем переменной первое число из последовательности и делаем в цикле до последнего AND
Последовательность от m до n
res = m
     for i in m..n
     res &= i
return res
А теперь интересное решение. Как вы видели, в последовательности, еденички выстроились в колонну. Напомню логику бинарного AND
    0101 (decimal 5)
AND 0011 (decimal 3)
  = 0001 (decimal 1)

Т.е. все по сути, все числа из последовательности уничножат друг друга. В AND будет 1 только когда в обоих числах бит равен 1. По сути, нам нужно найти этот префикс, где у 1-го и последнего числа будет этот префикс (он будет один, из условия задачи 0 <= m <= n <= 2147483647
Таким образом, примем, что число 32-х битное (из условия).  Поэтому O(1) будет циклом 0....32 (константное).
В этом цикле мы будем делать операцию битового сдвига вправо, пока не найдем этот общий префикс. Когда мы его нашли, нам нужно первое число, m сдвинуть уже влево на то количество бит, которое мы сдвигали вправо. Это делается для того, чтобы вернуть биты (хаха) обратно в число. Вопрос - а до куда в цикле идти, сколько сдвигать, как мы поймем, что нашли общий префикс. Сдвигается до того момента, пока число m меньше числа n. Так как когда они сравняются или будет n > m, значит у нас впереди только нули и сдвигать дальше не нужно. Пример в студию. Возьмем числа 9 и 12 (выше)
9   -- 0 0 0 0 1 0 0 1
12  -- 0 0 0 0 1 1 0 0
сдвигаем вправо, получаем 
4 -- 0 0 0 0 0 1 0 0
6 -- 0 0 0 0 0 1 1 0
m < n - идем дальше
2 -- 0 0 0 0 0 0 1 0
3 -- 0 0 0 0 0 0 1 1
m < n - идем дальше
1 -- 0 0 0 0 0 0 0 1
1 -- 0 0 0 0 0 0 0 1
Опа, m равно n. Отлично, самый дальний общий префикс найдет. Напомню, все, что было до этого будет нулями, когда мы будет делать 9&10, 10&11, 11&12 (не так конечно, но упрощенно, что я показывал выше)
Сдвиг равен 3. А теперь просто двигаем влево на 3 -- 0 0 0 0 1 0 0 0. 1 в четверном бите (третий индекс). 2 * 2 * 2 = 8.
Проверяем
9   -- 0 0 0 0 1 0 0 1
10  -- 0 0 0 0 1 0 1 0
11  -- 0 0 0 0 1 0 1 1
12  -- 0 0 0 0 1 1 0 0

Как видите, по пути от 9 к 12 если хоть у 1 числа есть ноль, там, где 1, все, там всегда будет 0. Т.е. получим 
res = 9&10 0 0 0 0 1 0 0 0
res&11 =     0 0 0 0 1 0 0 0
res&12 =    0 0 0 0 1 0 0 0
Готово)