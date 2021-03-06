package main

/*
 *    v6----- v5
 *   /|      /|
 *  v1------v0|
 *  | |     | |
 *  | v7----|-v4
 *  |/      |/
 *  v2------v3
 *
 * A cube has 6 sides and each side has 4 vertices, therefore, the total number
 * of vertices is 24 (6 sides * 4 verts), and 72 floats in the vertex array
 * since each vertex has 3 components (x,y,z) (= 24 * 3).
 */

/*
 * Vx, Vy, Vz -- Coordinate of the vertex;
 * Nx, Ny, Nz -- Coordinate of the normal;
 * U, V -- Coordinate of the texture;
 * R, G, B -- Color values.
 */

var cubeVertices = []float32{
	// Vx, Vy, Vz, Nx, Ny, Nz, U, V, R, G, B
	// Front face
	1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, // v0
	-1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, // v1
	-1, -1, 1, 0, 0, 1, 0, 1, 1, 0, 0, // v2
	1, -1, 1, 0, 0, 1, 1, 1, 1, 0, 1, // v3
	// Right face
	1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, // v0
	1, -1, 1, 1, 0, 0, 0, 1, 1, 0, 1, // v3
	1, -1, -1, 1, 0, 0, 1, 1, 0, 0, 1, // v4
	1, 1, -1, 1, 0, 0, 1, 0, 0, 1, 1, // v5
	// Top face
	1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, // v0
	1, 1, -1, 0, 1, 0, 1, 0, 0, 1, 1, // v5
	-1, 1, -1, 0, 1, 0, 0, 0, 0, 1, 0, // v6
	-1, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, // v1
	// Left face
	-1, 1, 1, -1, 0, 0, 1, 0, 1, 1, 0, // v1
	-1, 1, -1, -1, 0, 0, 0, 0, 0, 1, 0, // v6
	-1, -1, -1, -1, 0, 0, 0, 1, 0, 0, 0, // v7
	-1, -1, 1, -1, 0, 0, 1, 1, 1, 0, 0, // v2
	// Bottom face
	-1, -1, -1, 0, -1, 0, 0, 1, 0, 0, 0, // v7
	1, -1, -1, 0, -1, 0, 1, 1, 0, 0, 1, // v4
	1, -1, 1, 0, -1, 0, 1, 0, 1, 0, 1, // v3
	-1, -1, 1, 0, -1, 0, 0, 0, 1, 0, 0, // v2
	// Back face
	1, -1, -1, 0, 0, -1, 0, 1, 0, 0, 1, // v4
	-1, -1, -1, 0, 0, -1, 1, 1, 0, 0, 0, // v7
	-1, 1, -1, 0, 0, -1, 1, 0, 0, 1, 0, // v6
	1, 1, -1, 0, 0, -1, 0, 0, 0, 1, 1, // v5
}

var cubeIndices = []int32{
	// Front face
	0, 1, 2, // v0-v1-v2
	2, 3, 0, // v2-v3-v0
	// Right face
	4, 5, 6, // v0-v3-v4
	6, 7, 4, // v4-v5-v0
	// Top face
	8, 9, 10, // v0-v5-v6
	10, 11, 8, // v6-v1-v0
	// Left face
	12, 13, 14, // v1-v6-v7
	14, 15, 12, // v7-v2-v1
	// Bottom face
	16, 17, 18, // v7-v4-v3
	18, 19, 16, // v3-v2-v7
	// Back face
	20, 21, 22, // v4-v7-v6
	22, 23, 20, // v6-v5-v4
}
