# Laporan 03. Natural Cubic Spline Interpolation

## Data Input
- Titik 0: (-1, 2)
- Titik 1: (0, 2)
- Titik 2: (1, 9)
- Titik 3: (2, 4)
- Titik 4: (3, 4)

## Target Evaluasi
- Target 1: $x = -1/2 (≈ -0.5000)$
- Target 2: $x = 1/2 (≈ 0.5000)$
- Target 3: $x = 3/2 (≈ 1.5000)$
- Target 4: $x = 5/2 (≈ 2.5000)$

## Tahap 1: Pemetaan Jarak (h)
Rumus Umum: $h_i = x_{i+1} - x_i$

- $h_0 = 0 - -1 = 1$
- $h_1 = 1 - 0 = 1$
- $h_2 = 2 - 1 = 1$
- $h_3 = 3 - 2 = 1$

## Tahap 2: Syarat Batas
- $f''(x_0) = 0$
- $f''(x_4) = 0$

## Tahap 3: Sistem Tridiagonal untuk M_i
Dengan notasi $M_i = f''(x_i)$ dan syarat natural $M_0 = M_n = 0$, untuk tiap titik interior berlaku:

$h_{i-1}M_{i-1} + 2(h_{i-1}+h_i)M_i + h_iM_{i+1} = 6\left(\frac{y_{i+1}-y_i}{h_i} - \frac{y_i-y_{i-1}}{h_{i-1}}\right)$

Koefisien sistem selalu disusun dengan rumus umum agar tetap konsisten untuk jarak konstan maupun bervariasi.

- $i=1: 1M_0 + 4M_1 + 1M_2 = 42$
- $i=2: 1M_1 + 4M_2 + 1M_3 = -72$
- $i=3: 1M_2 + 4M_3 + 1M_4 = 30$

### Hasil Nilai $M_i = f''(x_i)$:
- $M_0 = f''(x_0) = 0$
- $M_1 = f''(x_1) = 237/14 (≈ 16.9286)$
- $M_2 = f''(x_2) = -180/7 (≈ -25.7143)$
- $M_3 = f''(x_3) = 195/14 (≈ 13.9286)$
- $M_4 = f''(x_4) = 0$

---
## Evaluasi untuk $x = -1/2 (≈ -0.5000)$

### Tahap 4: Pemilihan Segmen
- x = -1/2 (≈ -0.5000) berada pada interval [-1, 0], sehingga digunakan segmen 0.

### Tahap 5: Koefisien Segmen 0
- $a_0 = 2$
- $b_0 = -79/28 (≈ -2.8214)$
- $c_0 = 0$
- $d_0 = 79/28 (≈ 2.8214)$

### Tahap 6: Persamaan
$S(x) = 2 + (-79/28)(x - -1) + (0)(x - -1)^2 + (79/28)(x - -1)^3$

### Tahap 7: Evaluasi
- S(-1/2 (≈ -0.5000)) = 2 + (-79/28 (≈ -2.8214))(-1/2 (≈ -0.5000) - -1) + (0)(-1/2 (≈ -0.5000) - -1)^2 + (79/28 (≈ 2.8214))(-1/2 (≈ -0.5000) - -1)^3 = 211/224 (≈ 0.9420)

**Hasil Akhir: $y = 211/224 (≈ 0.9420)$**

---
## Evaluasi untuk $x = 1/2 (≈ 0.5000)$

### Tahap 4: Pemilihan Segmen
- x = 1/2 (≈ 0.5000) berada pada interval [0, 1], sehingga digunakan segmen 1.

### Tahap 5: Koefisien Segmen 1
- $a_1 = 2$
- $b_1 = 79/14 (≈ 5.6429)$
- $c_1 = 237/28 (≈ 8.4643)$
- $d_1 = -199/28 (≈ -7.1071)$

### Tahap 6: Persamaan
$S(x) = 2 + (79/14)(x - 0) + (237/28)(x - 0)^2 + (-199/28)(x - 0)^3$

### Tahap 7: Evaluasi
- S(1/2 (≈ 0.5000)) = 2 + (79/14 (≈ 5.6429))(1/2 (≈ 0.5000) - 0) + (237/28 (≈ 8.4643))(1/2 (≈ 0.5000) - 0)^2 + (-199/28 (≈ -7.1071))(1/2 (≈ 0.5000) - 0)^3 = 1355/224 (≈ 6.0491)

**Hasil Akhir: $y = 1355/224 (≈ 6.0491)$**

---
## Evaluasi untuk $x = 3/2 (≈ 1.5000)$

### Tahap 4: Pemilihan Segmen
- x = 3/2 (≈ 1.5000) berada pada interval [1, 2], sehingga digunakan segmen 2.

### Tahap 5: Koefisien Segmen 2
- $a_2 = 9$
- $b_2 = 5/4 (≈ 1.2500)$
- $c_2 = -90/7 (≈ -12.8571)$
- $d_2 = 185/28 (≈ 6.6071)$

### Tahap 6: Persamaan
$S(x) = 9 + (5/4)(x - 1) + (-90/7)(x - 1)^2 + (185/28)(x - 1)^3$

### Tahap 7: Evaluasi
- S(3/2 (≈ 1.5000)) = 9 + (5/4 (≈ 1.2500))(3/2 (≈ 1.5000) - 1) + (-90/7 (≈ -12.8571))(3/2 (≈ 1.5000) - 1)^2 + (185/28 (≈ 6.6071))(3/2 (≈ 1.5000) - 1)^3 = 1621/224 (≈ 7.2366)

**Hasil Akhir: $y = 1621/224 (≈ 7.2366)$**

---
## Evaluasi untuk $x = 5/2 (≈ 2.5000)$

### Tahap 4: Pemilihan Segmen
- x = 5/2 (≈ 2.5000) berada pada interval [2, 3], sehingga digunakan segmen 3.

### Tahap 5: Koefisien Segmen 3
- $a_3 = 4$
- $b_3 = -65/14 (≈ -4.6429)$
- $c_3 = 195/28 (≈ 6.9643)$
- $d_3 = -65/28 (≈ -2.3214)$

### Tahap 6: Persamaan
$S(x) = 4 + (-65/14)(x - 2) + (195/28)(x - 2)^2 + (-65/28)(x - 2)^3$

### Tahap 7: Evaluasi
- S(5/2 (≈ 2.5000)) = 4 + (-65/14 (≈ -4.6429))(5/2 (≈ 2.5000) - 2) + (195/28 (≈ 6.9643))(5/2 (≈ 2.5000) - 2)^2 + (-65/28 (≈ -2.3214))(5/2 (≈ 2.5000) - 2)^3 = 701/224 (≈ 3.1295)

**Hasil Akhir: $y = 701/224 (≈ 3.1295)$**
