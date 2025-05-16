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


# #!/usr/bin/env python3
# """
# Generate .mem files for compute_max_disp_tb with multiple test cases.

# Produces:
#   • left_image.mem   - TEST_CASES x (WIN x IMG_W) 8-bit hex values, one per line  
#   • right_image.mem  - same, but each window shifted by its random true_disp  
#   • exp_disp.mem     - TEST_CASES lines of 8-bit hex disparity values  
#   • exp_sad.mem      - TEST_CASES lines of “golden” SAD values (up to 16-bit)
# """
# import random

# # -----------------------------------------------------------------------------
# # PARAMETERS — must match your TB
# # -----------------------------------------------------------------------------
# WIN        = 15
# IMG_W      = 64
# MAX_DISP   = 64
# DATA_SIZE  = 8
# DATA_MAX   = (1 << DATA_SIZE) - 1
# TEST_CASES = 4

# # open all four files for writing
# with open("left_image.mem",  "w") as fL, \
#      open("right_image.mem", "w") as fR, \
#      open("exp_disp.mem",    "w") as fD, \
#      open("exp_sad.mem",     "w") as fS:

#     for tc in range(TEST_CASES):
#         # pick a random true disparity for this case
#         true_disp = random.randint(0, MAX_DISP - 1)

#         # build a WIN×IMG_W left window full of random pixels
#         left = [[random.randint(0, DATA_MAX) for _ in range(IMG_W)]
#                 for _ in range(WIN)]

#         # build the corresponding right window by shifting left→right
#         right = []
#         for r in range(WIN):
#             row = []
#             for c in range(IMG_W):
#                 if c >= true_disp:
#                     row.append(left[r][c - true_disp])
#                 else:
#                     row.append(random.randint(0, DATA_MAX))
#             right.append(row)

#         # flatten & write left window
#         for r in left:
#             for val in r:
#                 fL.write(f"{val:02X}\n")

#         # flatten & write right window
#         for r in right:
#             for val in r:
#                 fR.write(f"{val:02X}\n")

#         # write the expected disparity in HEX
#         fD.write(f"{true_disp:02X}\n")

#         # compute “golden” SAD for this shift
#         golden_sad = 0
#         for r in range(WIN):
#             for c in range(WIN):
#                 la = left[r][c]
#                 ra = right[r][c]  # already shifted
#                 golden_sad += abs(la - ra)
#         # write expected SAD in hex (up to 16 bits)
#         fS.write(f"{golden_sad:04X}\n")

# print(f"Generated {TEST_CASES} vectors in left_image.mem, right_image.mem,")
# print("and exp_disp.mem, exp_sad.mem")



# # Updated .mem generator with at least one “no-perfect-match” test case
# import random

# # -----------------------------------------------------------------------------
# # PARAMETERS — must match your TB
# # -----------------------------------------------------------------------------
# WIN        = 15
# IMG_W      = 64
# MAX_DISP   = 64
# DATA_SIZE  = 8
# DATA_MAX   = (1 << DATA_SIZE) - 1
# TEST_CASES = 4

# # pick one test case to have no perfect SAD=0 match
# no_zero_tc = random.randrange(TEST_CASES)

# # open all four files for writing
# with open("left_image_n.mem",  "w") as fL, \
#      open("right_image_n.mem", "w") as fR, \
#      open("exp_disp_n.mem",    "w") as fD, \
#      open("exp_sad_n.mem",     "w") as fS:

#     for tc in range(TEST_CASES):
#         # build a WIN×IMG_W left image of random pixels
#         left = [[random.randint(0, DATA_MAX) for _ in range(IMG_W)]
#                 for _ in range(WIN)]

#         # build right image
#         if tc == no_zero_tc:
#             # entirely random ⇒ no guaranteed SAD=0
#             right = [[random.randint(0, DATA_MAX) for _ in range(IMG_W)]
#                      for _ in range(WIN)]
#         else:
#             # shift by a random true_disp
#             true_disp = random.randint(0, MAX_DISP - 1)
#             right = []
#             for r in range(WIN):
#                 row = []
#                 for c in range(IMG_W):
#                     if c >= true_disp:
#                         row.append(left[r][c - true_disp])
#                     else:
#                         row.append(random.randint(0, DATA_MAX))
#                 right.append(row)

#         # flatten & write left image
#         for r in left:
#             for val in r:
#                 fL.write(f"{val:02X}\n")

#         # flatten & write right image
#         for r in right:
#             for val in r:
#                 fR.write(f"{val:02X}\n")

#         # compute “golden” best SAD and its disparity by exhaustive search
#         best_sad  = float('inf')
#         best_disp = 0
#         # only consider shifts where the window stays in bounds
#         max_valid_disp = IMG_W - WIN
#         for disp in range(min(MAX_DISP, max_valid_disp + 1)):
#             sad = 0
#             for r in range(WIN):
#                 for c in range(WIN):
#                     la = left[r][c]
#                     ra = right[r][c + disp]
#                     sad += abs(la - ra)
#             if sad < best_sad:
#                 best_sad  = sad
#                 best_disp = disp

#         # write expected disparity and SAD
#         fD.write(f"{best_disp:02X}\n")
#         fS.write(f"{best_sad:04X}\n")

# print(f"Generated {TEST_CASES} vectors (with no-perfect-match case at index {no_zero_tc})")


#!/usr/bin/env python3
import os
import random
import numpy as np
from PIL import Image

# -----------------------------------------------------------------------------
# PARAMETERS — must match your TB / C harness
# -----------------------------------------------------------------------------
WIN           = 15                   # SAD window size
IMG_W         = 128                  # patch width & height
MAX_DISP      = 64                   # max disparity
DATA_SIZE     = 8                    # bits/pixel
DATA_MAX      = (1 << DATA_SIZE) - 1
TEST_CASES    = 4

LEFT_IMG_FILE  = "L_00001.png"
RIGHT_IMG_FILE = "R_00001.png"
OUTPUT_DIR     = "mems"

# -----------------------------------------------------------------------------
# (No need to edit below here)
# -----------------------------------------------------------------------------
def compute_disparity_patch(left, right, max_disp, win):
    H, W = left.shape
    m = win // 2
    Lp = np.pad(left,  ((m, m), (m, m)), mode='constant', constant_values=0)
    Rp = np.pad(right, ((m, m), (m, m)), mode='constant', constant_values=0)

    best_cost = np.full((H, W), np.iinfo(np.int32).max, dtype=np.int32)
    best_disp = np.zeros((H, W), dtype=np.uint8)

    def integral_image(img):
        ii = img.cumsum(axis=0).cumsum(axis=1)
        out = np.zeros((ii.shape[0]+1, ii.shape[1]+1), dtype=np.int32)
        out[1:,1:] = ii
        return out

    for d in range(max_disp+1):
        if d>0:
            Rsh = np.zeros_like(Rp); Rsh[:,d:] = Rp[:,:-d]
        else:
            Rsh = Rp
        diff = np.abs(Lp.astype(np.int32) - Rsh.astype(np.int32))
        ii   = integral_image(diff)
        S = (
            ii[win:win+H, win:win+W]
          - ii[:H,      win:win+W]
          - ii[win:win+H, :W]
          + ii[:H,      :W]
        )
        mask = S < best_cost
        best_cost[mask] = S[mask]
        best_disp[mask] = d

    return best_disp

def write_mem(arr, path):
    with open(path, 'w') as f:
        for val in arr.flatten():
            f.write(f"{val:02X}\n")

def write_raw(arr, path):
    """Write a 2D uint8 array as raw bytes, row-major."""
    with open(path, 'wb') as f:
        f.write(arr.tobytes())

def main():
    # load full‐frame grayscale
    L = np.array(Image.open(LEFT_IMG_FILE).convert("L"), dtype=np.uint8)
    R = np.array(Image.open(RIGHT_IMG_FILE).convert("L"), dtype=np.uint8)
    H, W = L.shape

    os.makedirs(OUTPUT_DIR, exist_ok=True)

    left_mem  = open(f"{OUTPUT_DIR}/left_image.mem",  "w")
    right_mem = open(f"{OUTPUT_DIR}/right_image.mem", "w")
    disp_mem  = open(f"{OUTPUT_DIR}/exp_disp.mem",    "w")

    for tc in range(TEST_CASES):
        x0 = random.randint(0, W - IMG_W)
        y0 = random.randint(0, H - IMG_W)

        patch_L = L[y0:y0+IMG_W, x0:x0+IMG_W]
        patch_R = R[y0:y0+IMG_W, x0:x0+IMG_W]

        # dump raw for C harness
        write_raw(patch_L, f"{OUTPUT_DIR}/img_L_patch_{tc}.raw")
        write_raw(patch_R, f"{OUTPUT_DIR}/img_R_patch_{tc}.raw")

        # compute golden disparity patch
        disp_patch = compute_disparity_patch(patch_L, patch_R, MAX_DISP, WIN)

        # append to .mem files
        for p in patch_L.flatten():  left_mem.write(f"{p:02X}\n")
        for p in patch_R.flatten(): right_mem.write(f"{p:02X}\n")
        for d in disp_patch.flatten(): disp_mem.write(f"{d:02X}\n")

    left_mem.close()
    right_mem.close()
    disp_mem.close()

    print(f"Generated {TEST_CASES} patches of size {IMG_W}×{IMG_W} in `{OUTPUT_DIR}/`")
    print("Raw inputs saved as:")
    for tc in range(TEST_CASES):
        print(f"  mems/img_L_patch_{tc}.raw")
        print(f"  mems/img_R_patch_{tc}.raw")

if __name__ == "__main__":
    main()
