import os
import matplotlib.pyplot as plt
import numpy as np
import json
import codecs
import sys
import argparse
from json_numpy import object_hook



if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="k8s Insights"
    )

    parser.add_argument("--render-file", required=True)

    args = parser.parse_args()

    if args.render_file == "":
        print("--render-file should be non emtpy")
        exit(1)

    fname = args.render_file

    with codecs.open(fname,"rb" ,encoding="utf-8") as source:
        node = json.loads(source.read(),  object_hook=object_hook)

        print(type(node))
        print(node)

        node = np.fromiter(list(node), dtype=np.float64)


        median = np.median(node)
        print(median)

        plt.hist(node, bins=100, density=True)

        name = os.path.basename(fname)
        plt.title("Usage - " + name.split('.')[0])
        plt.xticks([0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1])
        plt.yticks([0,2,4,6,8,10,12,14,16,18,20])
        figure = plt.gcf()
        figure.savefig("out/" + name+ ".jpg")
        # print(node)