# Laporan 03. Natural Cubic Spline Interpolation

## Data Input
- Titik 0: (0, 2)
- Titik 1: (3, 3)
- Titik 2: (6, 1)
- Titik 3: (9, 0)
- Titik 4: (12, 7)

## Target Evaluasi
- Target 1: $x = 1$
- Target 2: $x = 5$
- Target 3: $x = 7$
- Target 4: $x = 11$

## Tahap 1: Pemetaan Jarak (h)
Rumus Umum: $h_i = x_{i+1} - x_i$

- $h_0 = 3 - 0 = 3$
- $h_1 = 6 - 3 = 3$
- $h_2 = 9 - 6 = 3$
- $h_3 = 12 - 9 = 3$

## Tahap 2: Syarat Batas
- $f''(x_0) = 0$
- $f''(x_4) = 0$

## Tahap 3: Sistem Tridiagonal untuk M_i
Dengan notasi $M_i = f''(x_i)$ dan syarat natural $M_0 = M_n = 0$, untuk tiap titik interior berlaku:

$h_{i-1}M_{i-1} + 2(h_{i-1}+h_i)M_i + h_iM_{i+1} = 6\left(\frac{y_{i+1}-y_i}{h_i} - \frac{y_i-y_{i-1}}{h_{i-1}}\right)$

Koefisien sistem selalu disusun dengan rumus umum agar tetap konsisten untuk jarak konstan maupun bervariasi.

- $i=1: 3M_0 + 12M_1 + 3M_2 = -6$
- $i=2: 3M_1 + 12M_2 + 3M_3 = 2$
- $i=3: 3M_2 + 12M_3 + 3M_4 = 16$

### Hasil Nilai $M_i = f''(x_i)$:
- $M_0 = f''(x_0) = 0$
- $M_1 = f''(x_1) = -41/84 (≈ -0.4881)$
- $M_2 = f''(x_2) = -1/21 (≈ -0.0476)$
- $M_3 = f''(x_3) = 113/84 (≈ 1.3452)$
- $M_4 = f''(x_4) = 0$

---
## Evaluasi untuk $x = 1$

### Tahap 4: Pemilihan Segmen
- x = 1 berada pada interval [0, 3], sehingga digunakan segmen 0.

### Tahap 5: Koefisien Segmen 0
- $a_0 = 2$
- $b_0 = 97/168 (≈ 0.5774)$
- $c_0 = 0$
- $d_0 = -33/1217 (≈ -0.0271)$

### Tahap 6: Persamaan
$S(x) = 2 + (97/168)(x - 0) + (0)(x - 0)^2 + (-33/1217)(x - 0)^3$

### Tahap 7: Evaluasi
- S(1) = 2 + (97/168 (≈ 0.5774))(1 - 0) + (0)(1 - 0)^2 + (-33/1217 (≈ -0.0271))(1 - 0)^3 = 482/189 (≈ 2.5503)

**Hasil Akhir: $y = 482/189 (≈ 2.5503)$**

---
## Evaluasi untuk $x = 5$

### Tahap 4: Pemilihan Segmen
- x = 5 berada pada interval [3, 6], sehingga digunakan segmen 1.

### Tahap 5: Koefisien Segmen 1
- $a_1 = 3$
- $b_1 = -13/84 (≈ -0.1548)$
- $c_1 = -41/168 (≈ -0.2440)$
- $d_1 = 37/1512 (≈ 0.0245)$

### Tahap 6: Persamaan
$S(x) = 3 + (-13/84)(x - 3) + (-41/168)(x - 3)^2 + (37/1512)(x - 3)^3$

### Tahap 7: Evaluasi
- S(5) = 3 + (-13/84 (≈ -0.1548))(5 - 3) + (-41/168 (≈ -0.2440))(5 - 3)^2 + (37/1512 (≈ 0.0245))(5 - 3)^3 = 361/189 (≈ 1.9101)

**Hasil Akhir: $y = 361/189 (≈ 1.9101)$**

---
## Evaluasi untuk $x = 7$

### Tahap 4: Pemilihan Segmen
- x = 7 berada pada interval [6, 9], sehingga digunakan segmen 2.

### Tahap 5: Koefisien Segmen 2
- $a_2 = 1$
- $b_2 = -23/24 (≈ -0.9583)$
- $c_2 = -1/42 (≈ -0.0238)$
- $d_2 = 13/168 (≈ 0.0774)$

### Tahap 6: Persamaan
$S(x) = 1 + (-23/24)(x - 6) + (-1/42)(x - 6)^2 + (13/168)(x - 6)^3$

### Tahap 7: Evaluasi
- S(7) = 1 + (-23/24 (≈ -0.9583))(7 - 6) + (-1/42 (≈ -0.0238))(7 - 6)^2 + (13/168 (≈ 0.0774))(7 - 6)^3 = 2/21 (≈ 0.0952)

**Hasil Akhir: $y = 2/21 (≈ 0.0952)$**

---
## Evaluasi untuk $x = 11$

### Tahap 4: Pemilihan Segmen
- x = 11 berada pada interval [9, 12], sehingga digunakan segmen 3.

### Tahap 5: Koefisien Segmen 3
- $a_3 = 0$
- $b_3 = 83/84 (≈ 0.9881)$
- $c_3 = 113/168 (≈ 0.6726)$
- $d_3 = -92/1231 (≈ -0.0747)$

### Tahap 6: Persamaan
$S(x) = 0 + (83/84)(x - 9) + (113/168)(x - 9)^2 + (-92/1231)(x - 9)^3$

### Tahap 7: Evaluasi
- S(11) = 0 + (83/84 (≈ 0.9881))(11 - 9) + (113/168 (≈ 0.6726))(11 - 9)^2 + (-92/1231 (≈ -0.0747))(11 - 9)^3 = 769/189 (≈ 4.0688)

**Hasil Akhir: $y = 769/189 (≈ 4.0688)$**
