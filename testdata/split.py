import os
import glob
import numpy as np
from PIL import Image
from skimage import color, transform


def split_interlaced_stereo_frames_2(dirname: str):
    """
    Input  - path to directory containing pgmyuv files extracted
             from a stereo video using mplayer.

    Output - left and right frames in the respective 'left' and
             'right' directories in the given input directory.
    """
    # Change to the specified directory
    os.chdir(dirname)

    # Get all pgmyuv files in the directory
    file_list = glob.glob("*.pgmyuv")

    # Create left and right directories if they don't exist
    os.makedirs("left", exist_ok=True)
    os.makedirs("right", exist_ok=True)

    for filename in file_list:
        print(f"processing {filename}")

        # Read the image and convert to floating point
        # In Python, we read the image and normalize to 0-1 range
        a = np.array(Image.open(filename)).astype(np.float32) / 255.0

        # Crop channels from the image
        Y = a[0:1080, 0:1920]
        U = a[1080:, 0:960]
        V = a[1080:, 960:]

        # Deinterlacing
        Y1 = Y[0::2, :]  # Every other row starting from 0
        Y2 = Y[1::2, :]  # Every other row starting from 1

        U1 = U[0::2, :]
        U2 = U[1::2, :]

        V1 = V[0::2, :]
        V2 = V[1::2, :]

        # Left, right splitting
        Y1left = Y1[:, 0:960]
        Y1right = Y1[:, 960:]

        Y2left = Y2[:, 0:960]
        Y2right = Y2[:, 960:]

        U1left = U1[:, 0:480]
        U1right = U1[:, 480:]

        U2left = U2[:, 0:480]
        U2right = U2[:, 480:]

        V1left = V1[:, 0:480]
        V1right = V1[:, 480:]

        V2left = V2[:, 0:480]
        V2right = V2[:, 480:]

        # Combine to RGB
        # Resize U and V channels (equivalent to MATLAB's imresize)
        U1left_resized = transform.resize(
            U1left, (U1left.shape[0], U1left.shape[1] * 2)
        )
        V1left_resized = transform.resize(
            V1left, (V1left.shape[0], V1left.shape[1] * 2)
        )

        U2left_resized = transform.resize(
            U2left, (U2left.shape[0], U2left.shape[1] * 2)
        )
        V2left_resized = transform.resize(
            V2left, (V2left.shape[0], V2left.shape[1] * 2)
        )

        U1right_resized = transform.resize(
            U1right, (U1right.shape[0], U1right.shape[1] * 2)
        )
        V1right_resized = transform.resize(
            V1right, (V1right.shape[0], V1right.shape[1] * 2)
        )

        U2right_resized = transform.resize(
            U2right, (U2right.shape[0], U2right.shape[1] * 2)
        )
        V2right_resized = transform.resize(
            V2right, (V2right.shape[0], V2right.shape[1] * 2)
        )

        # Stack the channels together
        ycbcr_1_left = np.stack((Y1left, U1left_resized, V1left_resized), axis=2)
        ycbcr_2_left = np.stack((Y2left, U2left_resized, V2left_resized), axis=2)
        ycbcr_1_right = np.stack((Y1right, U1right_resized, V1right_resized), axis=2)
        ycbcr_2_right = np.stack((Y2right, U2right_resized, V2right_resized), axis=2)

        # Convert from YCbCr to RGB
        rgb_1_left = color.ycbcr2rgb(ycbcr_1_left)
        rgb_2_left = color.ycbcr2rgb(ycbcr_2_left)
        rgb_1_right = color.ycbcr2rgb(ycbcr_1_right)
        rgb_2_right = color.ycbcr2rgb(ycbcr_2_right)

        # Save to individual files
        filestub = filename[:-7]  # Remove the .pgmyuv extension

        # Convert to uint8 and save as PPM
        Image.fromarray((rgb_1_left * 255).astype(np.uint8)).save(
            f"left/{filestub}_A.ppm"
        )
        Image.fromarray((rgb_2_left * 255).astype(np.uint8)).save(
            f"left/{filestub}_B.ppm"
        )
        Image.fromarray((rgb_1_right * 255).astype(np.uint8)).save(
            f"right/{filestub}_A.ppm"
        )
        Image.fromarray((rgb_2_right * 255).astype(np.uint8)).save(
            f"right/{filestub}_B.ppm"
        )

    # Change back to the parent directory
    os.chdir("..")


# Example usage:
# split_interlaced_stereo_frames_2('/path/to/directory')
