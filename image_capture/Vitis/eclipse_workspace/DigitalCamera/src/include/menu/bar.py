#!/usr/bin/env python3
import os
import glob
from PIL import Image
import re


def natural_sort_key(s):
    """Sort strings with embedded numbers naturally."""
    return [
        int(text) if text.isdigit() else text.lower() for text in re.split(r"(\d+)", s)
    ]


def ppm_to_gif(output_filename="animation.gif", frame_duration=100, loop=0):
    """
    Convert PPM files to a GIF animation.

    Parameters:
    - output_filename: Name of the output GIF file
    - frame_duration: Duration of each frame in milliseconds
    - loop: Number of times to loop the animation (0 = infinite)
    """
    # Get all PPM files in the current directory
    ppm_files = glob.glob("output_*.ppm")

    if not ppm_files:
        print("No PPM files found in the current directory.")
        return

    # Sort the files naturally so that they appear in the correct order
    ppm_files.sort(key=natural_sort_key)

    print(f"Found {len(ppm_files)} PPM files.")

    # Load the images
    images = []
    for ppm_file in ppm_files:
        try:
            img = Image.open(ppm_file)
            images.append(img)
            print(f"Loaded: {ppm_file}")
        except Exception as e:
            print(f"Error loading {ppm_file}: {e}")

    if not images:
        print("No images were successfully loaded.")
        return

    # Save as GIF
    try:
        images[0].save(
            output_filename,
            save_all=True,
            append_images=images[1:],
            duration=frame_duration,
            loop=loop,
            optimize=False,
        )
        print(f"GIF animation saved as {output_filename}")
    except Exception as e:
        print(f"Error saving GIF: {e}")


if __name__ == "__main__":
    # You can customize these parameters
    output_file = "animation.gif"
    frame_duration = 100  # milliseconds

    ppm_to_gif(output_file, frame_duration)
