#!/usr/bin/env python3
"""
@author:Harold
@file: displayTriangulation.py
@time: 27/09/2019
"""

import numpy as np
# This import registers the 3D projection, but is otherwise unused.
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt
from utils.utils import load_3d_pt_cloud_data_with_delimiter, set_axes_equal


fig = plt.figure(figsize=(20, 20))
ax = fig.add_subplot(111, projection="3d")

# plot projected points
normal_data = load_3d_pt_cloud_data_with_delimiter("ism_train_cat.txt", r"\s+")
X, Y, Z = (
    normal_data[:, 0],
    normal_data[:, 1],
    normal_data[:, 2],
)
ax.scatter(X, Y, Z, color="black")

# plot triangles
vertexes = load_3d_pt_cloud_data_with_delimiter("vertexID.txt", r"\s+")
for i in vertexes:
    v0, v1, v2 = int(i[0]), int(i[1]), int(i[2])
    if v0 >= 3400 or v1 >= 3400 or v2 >= 3400:
        continue
    edges = [(v0, v1), (v1, v2), (v2, v0)]
    for j, k in edges:
        ax.plot(normal_data[[j, k], 0], normal_data[[j, k], 1], normal_data[[j, k], 2], color='r')

set_axes_equal(ax)

ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
ax.view_init(45, 0)

plt.show()

