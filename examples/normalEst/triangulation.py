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
import math
'''
# https://github.com/mrzv/diode
# need first install CGAL
# `sudo apt-get install libcgal-dev`
# `sudo apt-get install libcgal-demo`
# `pip install --verbose git+https://github.com/mrzv/diode.git`
import diode
from dionysus import Filtration
'''


def load_3d_pt_cloud_data_with_delimiter(path_name: str, delimiter: str) -> np.array:
    return pd.read_csv(
        path_name, dtype=np.float32, delimiter=delimiter, header=None
    ).to_numpy()


# https://stackoverflow.com/questions/23073170/calculate-bounding-polygon-of-alpha-shape-from-the-delaunay-triangulation
def alpha_shape(points, alpha, only_outer=True):
    """
    Compute the alpha shape (concave hull) of a set of points.
    :param points: np.array of shape (n, 3) points.
    :param alpha: alpha value.
    :param only_outer: boolean value to specify if we keep only the outer border
    or also inner edges.
    :return: set of (i,j) pairs representing edges of the alpha-shape. (i,j) are
    the indices in the points array.
    """
    assert points.shape[0] > 3, "Need at least four points"

    def add_edge(edges, i, j):
        """
        Add a line between the i-th and j-th points,
        if not in the list already
        """
        if (i, j) in edges or (j, i) in edges:
            # already added
            if only_outer:
                # if both neighboring triangles are in shape, it's not a boundary edge
                if (j, i) in edges:
                    edges.remove((j, i))
            return
        edges.add((i, j))

    tri = Delaunay(points)
    edges = set()
    # Loop over triangles:
    # ia, ib, ic = indices of corner points of the triangle
    print(tri.vertices.shape)
    for ia, ib, ic, id in tri.vertices:
        pa = points[ia]
        pb = points[ib]
        pc = points[ic]
        pd = points[id]

        # Computing radius of triangle Circumsphere
        # http://mathworld.wolfram.com/Circumsphere.html

        pa2 = np.dot(pa, pa)
        pb2 = np.dot(pb, pb)
        pc2 = np.dot(pc, pc)
        pd2 = np.dot(pd, pd)

        a = np.linalg.det(np.array([np.append(pa, 1), np.append(pb, 1), np.append(pc, 1), np.append(pd, 1)]))

        Dx = np.linalg.det(np.array([np.array([pa2, pa[1], pa[2], 1]),
                                     np.array([pb2, pb[1], pb[2], 1]),
                                     np.array([pc2, pc[1], pc[2], 1]),
                                     np.array([pd2, pd[1], pd[2], 1])]))

        Dy = - np.linalg.det(np.array([np.array([pa2, pa[0], pa[2], 1]),
                                       np.array([pb2, pb[0], pb[2], 1]),
                                       np.array([pc2, pc[0], pc[2], 1]),
                                       np.array([pd2, pd[0], pd[2], 1])]))

        Dz = np.linalg.det(np.array([np.array([pa2, pa[0], pa[1], 1]),
                                     np.array([pb2, pb[0], pb[1], 1]),
                                     np.array([pc2, pc[0], pc[1], 1]),
                                     np.array([pd2, pd[0], pd[1], 1])]))

        c = np.linalg.det(np.array([np.array([pa2, pa[0], pa[1], pa[2]]),
                                    np.array([pb2, pb[0], pb[1], pb[2]]),
                                    np.array([pc2, pc[0], pc[1], pc[2]]),
                                    np.array([pd2, pd[0], pd[1], pd[2]])]))

        circum_r = math.sqrt(math.pow(Dx, 2) + math.pow(Dy, 2) + math.pow(Dz, 2) - 4 * a * c) / (2 * abs(a))
        if circum_r < alpha:
            add_edge(edges, ia, ib)
            add_edge(edges, ib, ic)
            add_edge(edges, ic, id)
            add_edge(edges, id, ia)
            add_edge(edges, ia, ic)
            add_edge(edges, ib, id)
    return edges


normal_data = load_3d_pt_cloud_data_with_delimiter("ism_train_cat_normal.txt", r"\s+")

original_data = normal_data[:, [0, 1, 2]]

hull = ConvexHull(original_data)

fig = plt.figure(figsize=(20, 20))
ax = fig.add_subplot(311, projection="3d")
for s in hull.simplices:
    # s = np.append(s, s[0])  # ycle back to the first coordinate
    ax.plot(original_data[s, 0], original_data[s, 1], original_data[s, 2], "r-")
ax.set_xlabel("x")
ax.set_ylabel("y")
ax.set_zlabel("z")
ax.view_init(45, 0)

tri = Delaunay(original_data)

bx = fig.add_subplot(312, projection="3d")
bx.plot_trisurf(original_data[:, 0], original_data[:, 1], original_data[:, 2], triangles=tri.simplices, cmap=plt.cm.Spectral)
bx.set_xlabel("x")
bx.set_ylabel("y")
bx.set_zlabel("z")
bx.view_init(45, 0)

'''
# seems useless...
# have to use CGAL directly...
simplices = diode.fill_alpha_shapes(original_data)
f = Filtration(simplices)
print(f)
alphashape = [s for s in f if s.data <= 0.1]
# print(alphashape[1000].data)
print(len(alphashape))
cx = fig.add_subplot(313, projection="3d")
cx.plot_trisurf(original_data[:, 0], original_data[:, 1], original_data[:, 2], triangles=alphashape, cmap=plt.cm.Spectral)
cx.set_xlabel("x")
cx.set_ylabel("y")
cx.set_zlabel("z")
cx.view_init(45, 0)
'''

edges = alpha_shape(original_data, 5)
print(len(edges))
cx = fig.add_subplot(313, projection="3d")
# cx.scatter(original_data[:, 0], original_data[:, 1], original_data[:, 2], '.')
for i, j in edges:
    cx.plot(original_data[[i, j], 0], original_data[[i, j], 1], original_data[[i, j], 2], color='r')
cx.set_xlabel("x")
cx.set_ylabel("y")
cx.set_zlabel("z")
cx.view_init(45, 0)

plt.show()
