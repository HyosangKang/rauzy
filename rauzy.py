import numpy as np
import numpy.linalg as ln
import matplotlib.pyplot as plt

class Rauzy:
    def __init__(self, subs, n):
        pts = [[] for i in range(3)]
        w = [0]
        for i in range(n):
            wt = []
            for j in range(len(w)):
                wt.extend(subs[w[j]])
            w = wt
        self.ev = np.array([0.0, 0.0, 0.0])
        for i in w:
            self.ev[i] += 1
            pts[i].append(np.array(self.ev))
        self.ev /= ln.norm(self.ev)
        a = np.array([self.ev, [1.0, 0.0, 0.0], [0.0, 1.0, 0.0]]).transpose()
        q, _ = ln.qr(a)
        q = q.transpose()
        e = [q[1], q[2]]
        self.cds = [[[] for _ in range(2)] for _ in range(3)]
        for i in range(3):
            for pt in pts[i]:
                for j in range(2):
                    self.cds[i][j].append(np.dot(pt, e[j]))
    def draw(self):
        for i in range(3):
            plt.scatter(self.cds[i][0], self.cds[i][1], s=1)
        plt.show()
    def eigenvector(self):
        return self.ev

if __name__ == '__main__':
    r = Rauzy([[0, 1], [0, 2], [0]], 20)
    print(r.eigenvector())
    r.draw()