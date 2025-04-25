import numpy as np

## Only used for display and image import
import cv2
import matplotlib.pyplot as plt

def cv2_Disparity(left_source, right_source, block_size, max_disparity):
    stereo = cv2.StereoBM_create(numDisparities=max_disparity, blockSize=block_size)
    disparity = stereo.compute(left_source, right_source)
    # min = disparity.min()
    # max = disparity.max()
    # disparity = np.uint8(255 * (disparity - min) / (max - min))
    return disparity


def SAD(left_source, right_source, block_size, max_disparity):
    h, w = left_source.shape
    disparity_map = np.zeros((h, w), dtype=int)

    half_block = block_size // 2

    for y in range(half_block, h - half_block):
        for x in range(half_block, w - half_block):
            best_offset = 0
            min_sad = float('inf')

            for d in range(max_disparity):
                x_shifted = x - d
                if x_shifted < half_block:
                    break

                left_block = left_source[y - half_block:y + half_block + 1, x - half_block:x + half_block + 1]
                right_block = right_source[y - half_block:y + half_block + 1, x_shifted - half_block:x_shifted + half_block + 1]

                sad = np.sum(np.abs(left_block - right_block))

                if sad < min_sad:
                    min_sad = sad
                    best_offset = d

                disparity_map[y, x] = best_offset

    return disparity_map.astype(np.uint8)

img_L = cv2.imread("L_00001.png", cv2.IMREAD_GRAYSCALE)
img_R = cv2.imread("R_00001.png", cv2.IMREAD_GRAYSCALE)

disparity_map_cv2 = cv2_Disparity(img_L, img_R, 5, 80)
disparity_map = SAD(img_L, img_R, 5, 80)

# disparity_map_filtered = cv2.medianBlur(disparity_map.astype(np.uint8), 5)


plt.figure(num=1, figsize=(10,5))
plt.subplot(1,2,1)
plt.imshow(img_L, 'gray')
plt.title("Original")
plt.subplot(1,2,2)
plt.imshow(disparity_map_cv2, 'gray')
plt.title("disparity_map_cv2")

plt.figure(num=2, figsize=(10,5))
plt.subplot(1,2,1)
plt.imshow(img_L, 'gray')
plt.title("Original")
plt.subplot(1,2,2)
plt.imshow(disparity_map, 'gray')
plt.title("disparity_map")
plt.show()

# # Display the result
# plt.imshow(disparity_map, cmap='gray')
# plt.colorbar()
# plt.title("Disparity Map")
# plt.show()


# def png2arr(img2convert):
#     arr = np.array(img2convert)
#     return arr