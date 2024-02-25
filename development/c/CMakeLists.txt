cmake_minimum_required(VERSION 3.26)
project(c C)

set(CMAKE_CXX_FLAGS_RELEASE "-DNDEBUG -Wall -O2")
set(CMAKE_CXX_FLAGS_DEBUG "-DNDEBUG -Wall -O2")
set(CMAKE_C_STANDARD 11)

add_executable(c
        bresenham.c
        lodepng.c)
