#!/usr/bin/env python3
"""
@author:Harold
@file: normal.py
@time: 19/09/2019
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


def load_3d_pt_cloud_data_with_delimiter(path_name: str, delimiter: str) -> np.array:
    return pd.read_csv(path_name, dtype=np.float32, delimiter=delimiter, header=None).to_numpy()


fig = plt.figure(figsize=(20, 10))

normal_data = load_3d_pt_cloud_data_with_delimiter('ism_train_cat_normal.txt', r"\s+")

X, Y, Z, U, V, W = normal_data[:, 0], normal_data[:, 1], normal_data[:, 2], normal_data[:, 3], normal_data[:, 4], normal_data[:, 5]

ax = fig.add_subplot(111, projection='3d')
ax.scatter(X, Y, Z, color='red')
ax.quiver(X, Y, Z, U, V, W, pivot='tail', length=3)
ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
ax.view_init(45, 0)
# plt.savefig("normal.svg")
plt.show()
