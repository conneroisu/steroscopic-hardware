# # save as convert_to_raw.py
# import cv2
# import numpy as np

# def convert_to_grayscale_raw(input_path, output_path):
#     img = cv2.imread(input_path, cv2.IMREAD_GRAYSCALE)
#     if img is None:
#         raise ValueError(f"Failed to read {input_path}")
#     img.tofile(output_path)

# convert_to_grayscale_raw("C:\\Users\\jaxie\\Desktop\\cpre488\\steroscopic-hardware\\L_00001.png", "img_L.raw")
# convert_to_grayscale_raw("C:\\Users\\jaxie\\Desktop\\cpre488\\steroscopic-hardware\\R_00001.png", "img_R.raw")

import numpy as np
import matplotlib.pyplot as plt

img = np.fromfile("C:\\Users\\jaxie\\Desktop\\cpre488\\steroscopic-hardware\\disparity.pgm", dtype=np.uint8, offset=15).reshape((480, 640))
plt.imshow(img, cmap='gray')
plt.colorbar()
plt.title("Disparity Map")
plt.show()