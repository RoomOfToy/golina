#!/usr/bin/env python3
"""
@author:Harold
@file: mesh.py
@time: 20/09/2019
"""

# This import registers the 3D projection, but is otherwise unused.
from mpl_toolkits.mplot3d import Axes3D
import matplotlib

# display interactively
# `sudo apt-get install python3-tk` to install tkinter first
matplotlib.use("TkAgg")
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
from scipy.spatial import ConvexHull, Delaunay


def load_3d_pt_cloud_data_with_delimiter(path_name: str, delimiter: str) -> np.array:
    return pd.read_csv(
        path_name, dtype=np.float32, delimiter=delimiter, header=None
    ).to_numpy()


normal_data = load_3d_pt_cloud_data_with_delimiter("ism_train_cat_normal.txt", r"\s+")

original_data = normal_data[:, [0, 1, 2]]

hull = ConvexHull(original_data)

fig = plt.figure(figsize=(20, 20))
ax = fig.add_subplot(211, projection="3d")
for s in hull.simplices:
    # s = np.append(s, s[0])  # ycle back to the first coordinate
    ax.plot(original_data[s, 0], original_data[s, 1], original_data[s, 2], "r-")
ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
ax.view_init(45, 0)

tri = Delaunay(original_data)

bx = fig.add_subplot(212, projection="3d")
bx.plot_trisurf(original_data[:, 0], original_data[:, 1], original_data[:, 2], triangles=tri.simplices, cmap=plt.cm.Spectral)
bx.set_xlabel("x")
bx.set_ylabel("y")
bx.set_zlabel("z")
bx.view_init(45, 0)

plt.show()
