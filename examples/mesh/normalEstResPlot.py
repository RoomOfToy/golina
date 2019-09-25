#!/usr/bin/env python3
"""
@author:Harold
@file: mesh.py
@time: 18/09/2019
"""

# This import registers the 3D projection, but is otherwise unused.
from mpl_toolkits.mplot3d import Axes3D

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


def load_3d_pt_cloud_data_with_delimiter(path_name: str, delimiter: str) -> np.array:
    return pd.read_csv(
        path_name, dtype=np.float32, delimiter=delimiter, header=None
    ).to_numpy()


fig = plt.figure(figsize=(20, 10))

# estimated normal data
normal_data = load_3d_pt_cloud_data_with_delimiter("ism_train_cat_normal.txt", r"\s+")
X, Y, Z, U, V, W = (
    normal_data[:, 0],
    normal_data[:, 1],
    normal_data[:, 2],
    normal_data[:, 3],
    normal_data[:, 4],
    normal_data[:, 5],
)

bx = fig.add_subplot(211, projection="3d")
bx.scatter(X, Y, Z, color="red")
bx.quiver(X, Y, Z, U, V, W)
bx.set_xlabel("x")
bx.set_ylabel("y")
bx.set_zlabel("z")
bx.view_init(45, 0)

# original_data
original_data = load_3d_pt_cloud_data_with_delimiter("ism_train_cat.txt", r"\s+")
x, y, z = original_data[:, 0], original_data[:, 1], original_data[:, 2]

cx = fig.add_subplot(212, projection="3d")
cx.scatter(x, y, z, color="black")
cx.set_xlabel("x")
cx.set_ylabel("y")
cx.set_zlabel("z")
cx.view_init(45, 0)

plt.show()
