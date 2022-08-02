import numpy as np
import numpy.linalg as ln
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D

class Rauzy:
    subs = [] # substitution rule
    dim = 0 # dimension of subsitution
    cds = [] # projected coordinates
    word = [] # morphic word
    
    def __init__(self, subs):
        self.subs = subs
        self.dim = len(subs)    
        self.word = [0]

    def draw(self, sz=5):
        if self.dim == 3:
            for i in range(self.dim):
                plt.scatter(self.cds[i][0], self.cds[i][1], s=sz)
        if self.dim == 4:
            ax = plt.axes(projection='3d') 
            for i in range(self.dim):
                ax.scatter3D(self.cds[i][0], self.cds[i][1], self.cds[i][2], s=sz)
        plt.show()

    def morph(self):
        w = [] 
        for c in self.word:
            w.extend(self.subs[c])
        self.word = w

    def eigenvector(self):
        ev = [0 for _ in range(self.dim)]
        for i in self.word:
            ev[i] += 1
        ev = np.array(ev)
        self.ev = ev / ln.norm(ev)

    def project(self):
        m = [self.ev]
        for i in range(self.dim-1):
            v = [0.0 for _ in range(self.dim)]
            v[i] = 1.0
            m.append(v)
        a = np.array(m).transpose()
        q, _ = ln.qr(a)
        q = q.transpose()
        e = []
        for i in range(1, self.dim):
            e.append(q[i])
        self.cds = [[[] for _ in range(self.dim-1)] for _ in range(self.dim)]
        v = np.array([0 for _ in range(self.dim)])
        for c in self.word:
            v[c] += 1
            for i in range(self.dim-1):
                self.cds[c][i].append(np.dot(v, e[i]))

    def run(self, n):
        for i in range(n):
            self.morph()
        self.eigenvector()
        self.project()

if __name__ == '__main__':
    # r = Rauzy([[0, 1], [0, 2], [0]])
    r = Rauzy([[0, 1], [0, 2], [0, 3], [0]])
    r.run(15)
    r.draw(sz=5)