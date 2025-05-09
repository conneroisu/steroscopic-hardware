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


# ChatGPT generated for .mem files 
import random

def generate_mem(filename, size, pattern='increment', base=0):
    values = []

    for i in range(size):
        if pattern == 'increment':
            val = (base + i) % 256
        elif pattern == 'constant':
            val = base
        elif pattern == 'random':
            val = random.randint(0, 255)
        elif pattern == 'checker':
            val = 255 if (i % 2 == 0) else 0
        else:
            raise ValueError("Unknown pattern")

        values.append(f"{val:02X}")

    # Write hex values to file
    with open(filename, 'w') as f:
        for v in values:
            f.write(v + '\n')

    print(f"[+] Generated {size} values in {filename} using '{pattern}' pattern.")

# Example usage
if __name__ == "__main__":
    WIN = 15
    WIN_SIZE = WIN * WIN

    generate_mem("winL.mem", WIN_SIZE, pattern='increment', base=0)
    generate_mem("winR.mem", WIN_SIZE, pattern='random')