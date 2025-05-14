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

# import numpy as np
# import matplotlib.pyplot as plt

# img = np.fromfile("C:\\Users\\jaxie\\Desktop\\cpre488\\steroscopic-hardware\\disparity.pgm", dtype=np.uint8, offset=15).reshape((480, 640))
# plt.imshow(img, cmap='gray')
# plt.colorbar()
# plt.title("Disparity Map")
# plt.show()


# #!/usr/bin/env python3
# """
# Generate .mem files for compute_max_disp_tb with TEST_CASES=1.

# Produces:
#   • left_image.mem   – WIN × IMG_W values of left window
#   • right_image.mem  – WIN × IMG_W values of right window,
#                        shifted by a random true_disp (0 ≤ true_disp < MAX_DISP)
#   • exp_disp.mem     – single-line containing the true_disp

# All values are 8-bit hex, one per line, ready for $readmemh().
# """
# import random

# # parameters (must match your TB)
# WIN       = 15
# IMG_W     = 64
# MAX_DISP  = 64
# DATA_MAX  = (1 << 8) - 1  # 8-bit pixels
# TEST_CASES = 1             # fixed to 1

# # pick a random “true” disparity
# true_disp = random.randint(0, MAX_DISP - 1)

# # build left image block
# left = [
#     [random.randint(0, DATA_MAX) for _ in range(IMG_W)]
#     for _ in range(WIN)
# ]

# # build right image block by shifting left → right by true_disp
# right = []
# for r in range(WIN):
#     row = []
#     for c in range(IMG_W):
#         if c >= true_disp:
#             # alignment: right pixel at c came from left at c - true_disp
#             row.append(left[r][c - true_disp])
#         else:
#             # fill the “underflow” with random noise
#             row.append(random.randint(0, DATA_MAX))
#     right.append(row)

# # helper to flatten and write a .mem file
# def write_mem(filename, data_2d):
#     with open(filename, "w") as f:
#         for row in data_2d:
#             for val in row:
#                 f.write(f"{val:02X}\n")

# # write left_image.mem and right_image.mem
# write_mem("left_image.mem",  left)
# write_mem("right_image.mem", right)

# # write exp_disp.mem (one line)
# with open("exp_disp.mem", "w") as f:
#     f.write(f"{true_disp}\n")

# print(f"Generated left_image.mem, right_image.mem (each {WIN*IMG_W} lines) ")
# print(f"and exp_disp.mem with true_disp = {true_disp}")


#!/usr/bin/env python3
"""
Generate .mem files for compute_max_disp_tb with multiple test cases.

Produces:
  • left_image.mem   – TEST_CASES×(WIN × IMG_W) 8-bit hex values, one per line
  • right_image.mem  – same, but each window shifted by its random true_disp
  • exp_disp.mem     – TEST_CASES lines of 8-bit hex disparity values

All values are zero-padded hex (00–FF), ready for $readmemh().
"""
import random

# -----------------------------------------------------------------------------
# PARAMETERS — must match your TB
# -----------------------------------------------------------------------------
WIN         = 15
IMG_W       = 64
MAX_DISP    = 64
DATA_SIZE   = 8
DATA_MAX    = (1 << DATA_SIZE) - 1
TEST_CASES  = 4    # change to how many vectors you want

# open all three files for writing
with open("left_image.mem",  "w") as fL, \
     open("right_image.mem", "w") as fR, \
     open("exp_disp.mem",    "w") as fE:

    for tc in range(TEST_CASES):
        # pick a random true disparity for this case
        true_disp = random.randint(0, MAX_DISP - 1)

        # build a WIN×IMG_W left window full of random pixels
        left = [[random.randint(0, DATA_MAX) for _ in range(IMG_W)]
                for _ in range(WIN)]

        # build the corresponding right window by shifting by true_disp
        right = []
        for r in range(WIN):
            row = []
            for c in range(IMG_W):
                if c >= true_disp:
                    # pixel came from left[r][c - true_disp]
                    row.append(left[r][c - true_disp])
                else:
                    # underflow region → random noise
                    row.append(random.randint(0, DATA_MAX))
            right.append(row)

        # flatten & write left window
        for r in left:
            for val in r:
                fL.write(f"{val:02X}\n")

        # flatten & write right window
        for r in right:
            for val in r:
                fR.write(f"{val:02X}\n")

        # write the expected disparity in HEX
        fE.write(f"{true_disp:02X}\n")

print(f"Generated {TEST_CASES} test vectors in left_image.mem, right_image.mem, exp_disp.mem")
