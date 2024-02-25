#include "lodepng.h"

#include <stdio.h>
#include <time.h>
#include <stdlib.h>

void encode(const char* filename, const unsigned char* p, unsigned w, unsigned h) {
    unsigned error = lodepng_encode32_file(filename, p, w, h);
    if(error) printf("error %u: %s\n", error, lodepng_error_text(error));
}

int gcd(int a, int b ) {
    if (a == 0) {
        return b;
    }
    while (b != 0) {
        int h = a % b;
        a = b;
        b = h;
    }
    return a;
}

void bresenham2(const int x1, const int y1 , const int x2, const int y2, unsigned int * p, const unsigned int c, const int w, const int h) {
    const int dx = x2 - x1;
	const int dy = y2 - y1;

	if (dx == 0) {
		return;
	}
	if (dy == 0) {
		int i = w*y1 + x1;
		int end = w*y2 + x2;
		while (i <= end) {
			p[i] = c;
			i++;
			p[end] = c;
			end--;
		}
		return;
	}
	if (dx == dy) {
		int i = w*y1 + x1;
        int end = w*y2 + x2;
        int off = 1 + w;
		while (i <= end) {
			p[i] = c;
			i += off;
			p[end] = c;
			end -= off;
		}
		return;
	}

    const int iMod[2] = {1, w + 1};

	int i = w*y1 + x1;

	const int g = gcd(dx, dy);
	const int goff = dx + dy*w;

	const int ndy2 = -(dy << 1); // negated
	const int dx2 = dx << 1;
	int e = ndy2 + dx; // e is negated
	p[i] = c;

	const int gpo = goff / g;

	int end = i + gpo;

	for (int j = i + gpo; j <= i+goff; j += gpo) {
		p[j] = c;
	}

	p[end] = c;

	int b = (int)((unsigned int)(e) >> 31);
    int z = (int)(~((unsigned int)(-e)) >> 31);
	e += dx2*b + ndy2;
	i += iMod[b];
	end -= iMod[z];

	if (g > 1) {
		while (i <= end) {
            b = (int)((unsigned int)(e) >> 31);
            z = (int)(~((unsigned int)(-e)) >> 31);
			p[i] = c;
			p[end] = c;
			const int off = end - i;
            int j = i + gpo;
            const int m = i + goff;
			while (j < m) {
				p[j] = c;
				p[j+off] = c;
				j += gpo;
			}
			e += ndy2;
			i += iMod[b];
			end -= iMod[z];
			e += dx2 * b;
		}
	} else {
		while (i <= end) {
			b = (int)((unsigned int)(e) >> 31);
			z = (int)(~((unsigned int)(-e)) >> 31);
			p[i] = c;
			p[end] = c;

			e += ndy2;
			i += iMod[b];
			end -= iMod[z];
			e += dx2 * b;
		}
	}
}

void bresenham(int x1, int y1 , int x2, int y2, unsigned int * p, unsigned int c,  int w,  int h) {

    int iMod[2] = {1, w + 1};

    int i = w*y1 + x1;
    p[i] = c;
    int end = w*y2 + x2;
    p[end] = c;

    int dx = x2 - x1;
    int ndy2 = -((y2 - y1) << 1); //negated
    int dx2 = dx << 1;
    int e = ndy2 + dx; //e is negated   

    int b = (int)(((unsigned int)e) >> 31);
    e += dx2*b + ndy2;
    i += iMod[b];

    while (i < end) {
        b = (int)(((unsigned int)e) >> 31);
        p[i] = c;
        e += ndy2;
        i += iMod[b];
        e += dx2 * b;
    }
}

int randN(int x) {
    return rand() % x;
}

int min(int x, int y) {
    if (x < y) {
        return x;
    }
    return y;
}

int max(int x,int y) {
    if (x > y) {
        return x;
    }
    return y;
}

void benchmarkSame(int n, int res, int num_raster, void (**rasterizer)(int , int  , int , int , unsigned int * , unsigned int ,  int ,  int)) {
    unsigned int pixel[res * res];
    for (int i = 0; i < n; i++) {
        unsigned int c = ((unsigned int)rand()) | (0b11111111 << 24);
        int r1 = randN(res);
        int r2 = randN(res);
        int r3 = randN(res);
        int r4 = randN(res);
        int x1 = min(r1, r2);
        int x2 = max(r1, r2);
        int y1 = min(r3, r4);
        int y2 = max(r3, r4);
        int dx = x2 - x1;
        int dy = y2 - y1;
        if (dx < dy) {
            int tmp = y1;
            y1 = x1;
            x1 = tmp;

            tmp = y2;
            y2 = x2;
            x2 = tmp;
        }
        for (int index = 0; index < num_raster; index++) {
            (*rasterizer[index])(x1, y1, x2, y2, pixel, c, res , res);
        }
    }

    encode("out.png", (const unsigned char *) pixel, res, res);
}

int main() {
    srand(time(NULL));

    void (*rasterizer [2])(int , int  , int , int , unsigned int * , unsigned int ,  int ,  int);
    rasterizer[0] = bresenham;
    rasterizer[1] = bresenham2;

    benchmarkSame(1<<25, 1<<8, 2,rasterizer);
    benchmarkSame(1<<16, 1<<14, 2,rasterizer);

    return 0;
}