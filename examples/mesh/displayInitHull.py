#!/usr/bin/env python3
"""
@author:Harold
@file: displayInitHull.py
@time: 27/09/2019
"""

'''
Triangle ID: 1
        Vertex[0]: {0.930852, 0.359951, 0.062847}
        Vertex[1]: {-0.007683, -0.017061,  0.999825}
        Vertex[2]: {-0.149304,  0.988556,  0.021570}
        Neighbors[0] ID: 2      Neighbors[0] ID: 5      Neighbors[0] ID: 4

Triangle ID: 2
        Vertex[0]: {0.930852, 0.359951, 0.062847}
        Vertex[1]: {-0.037406, -0.998846,  0.030123}
        Vertex[2]: {-0.007683, -0.017061,  0.999825}
        Neighbors[0] ID: 3      Neighbors[0] ID: 6      Neighbors[0] ID: 1

Triangle ID: 3
        Vertex[0]: {0.930852, 0.359951, 0.062847}
        Vertex[1]: {-0.000000, -0.000000, -1.000000}
        Vertex[2]: {-0.037406, -0.998846,  0.030123}
        Neighbors[0] ID: 4      Neighbors[0] ID: 7      Neighbors[0] ID: 2

Triangle ID: 4
        Vertex[0]: {0.930852, 0.359951, 0.062847}
        Vertex[1]: {-0.149304,  0.988556,  0.021570}
        Vertex[2]: {-0.000000, -0.000000, -1.000000}
        Neighbors[0] ID: 1      Neighbors[0] ID: 8      Neighbors[0] ID: 3

Triangle ID: 5
        Vertex[0]: {-1.000000, -0.000000, -0.000000}
        Vertex[1]: {-0.149304,  0.988556,  0.021570}
        Vertex[2]: {-0.007683, -0.017061,  0.999825}
        Neighbors[0] ID: 8      Neighbors[0] ID: 1      Neighbors[0] ID: 6

Triangle ID: 6
        Vertex[0]: {-1.000000, -0.000000, -0.000000}
        Vertex[1]: {-0.007683, -0.017061,  0.999825}
        Vertex[2]: {-0.037406, -0.998846,  0.030123}
        Neighbors[0] ID: 5      Neighbors[0] ID: 2      Neighbors[0] ID: 7

Triangle ID: 7
        Vertex[0]: {-1.000000, -0.000000, -0.000000}
        Vertex[1]: {-0.037406, -0.998846,  0.030123}
        Vertex[2]: {-0.000000, -0.000000, -1.000000}
        Neighbors[0] ID: 6      Neighbors[0] ID: 3      Neighbors[0] ID: 8

Triangle ID: 8
        Vertex[0]: {-1.000000, -0.000000, -0.000000}
        Vertex[1]: {-0.000000, -0.000000, -1.000000}
        Vertex[2]: {-0.149304,  0.988556,  0.021570}
        Neighbors[0] ID: 7      Neighbors[0] ID: 4      Neighbors[0] ID: 5
'''

import numpy as np
# This import registers the 3D projection, but is otherwise unused.
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt
from utils.utils import load_3d_pt_cloud_data_with_delimiter, set_axes_equal


vertexes = np.array([
    [ 0.930852,  0.359951,  0.062847],
    [-0.007683, -0.017061,  0.999825],
    [-0.149304,  0.988556,  0.021570],
    [-0.037406, -0.998846,  0.030123],
    [-0.000000, -0.000000, -1.000000],
    [-1.000000, -0.000000, -0.000000]
])

fig = plt.figure(figsize=(20, 20))
ax = fig.add_subplot(111, projection="3d")

# plot init hull vertexes
X, Y, Z = vertexes[:, 0], vertexes[:, 1], vertexes[:, 2]
ax.scatter(X, Y, Z, 'r')
ax.plot_trisurf(vertexes[:, 0], vertexes[:, 1], vertexes[:, 2])

# plot projected points
normal_data = load_3d_pt_cloud_data_with_delimiter("projected_points.txt", r"\s+")
X, Y, Z = (
    normal_data[:, 0],
    normal_data[:, 1],
    normal_data[:, 2],
)
ax.scatter(X, Y, Z, color="red")

# plot unit sphere
u, v = np.mgrid[0:2*np.pi:200j, 0:np.pi:100j]
x = np.cos(u)*np.sin(v)
y = np.sin(u)*np.sin(v)
z = np.cos(v)
ax.plot_wireframe(x, y, z, color="g", alpha=0.5)

set_axes_equal(ax)

ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
ax.view_init(45, 0)

plt.show()
