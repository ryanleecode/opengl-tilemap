cmake_minimum_required(VERSION 3.16)
project(opengl_tilemap C)

add_subdirectory(glad)

set(CMAKE_C_STANDARD 11)
add_compile_options(-Wall -Wextra)
add_executable(opengl_tilemap main.c)
target_link_libraries(opengl_tilemap PRIVATE glad)
target_link_libraries(opengl_tilemap PRIVATE "-lglfw3 -lGL -lX11 -lpthread -lXrandr -lXi -ldl -lm")
