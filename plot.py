import matplotlib.pyplot as plt
import pandas as pd
import numpy as np


# Load the second CSV file
file_path_2 = './mrg_sequential_vs_parallel.csv'

# Read the second CSV file
data_2 = pd.read_csv(file_path_2)

# Load the uploaded CSV file to inspect its content
file_path = './mrg_basic_vs_parallel.csv'

# Read the CSV file
data = pd.read_csv(file_path)

# Data from the first file
sizes_1 = data['SIZE']
basic_times = data['Basic']
p1_times_1 = data['p=1']
p3_times_1 = data['p=3']
p6_times_1 = data['p=6']

# Data from the second file
sizes_2 = data_2['SIZE']
sequential_times = data_2['Sequential']
p1_times_2 = data_2['p=1']
p3_times_2 = data_2['p=3']
p6_times_2 = data_2['p=6']

# Plotting the first comparison (Basic vs Parallel)
plt.figure(figsize=(10, 6))
plt.plot(sizes_1, basic_times, label="Basic Merge", marker='o')
plt.plot(sizes_1, p1_times_1, label="Parallel Merge (p=1)", marker='o')
plt.plot(sizes_1, p3_times_1, label="Parallel Merge (p=3)", marker='o')
plt.plot(sizes_1, p6_times_1, label="Parallel Merge (p=6)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Basic vs Parallel Merge (Different Thread Levels)')
plt.legend()
plt.grid(True)
plt.savefig('plot_basic_vs_paraller_merge.png')

# Plotting the second comparison (Sequential vs Parallel)
plt.figure(figsize=(10, 6))
plt.plot(sizes_2, sequential_times, label="Sequential Merge", marker='o')
plt.plot(sizes_2, p1_times_2, label="Parallel Merge (p=1)", marker='o')
plt.plot(sizes_2, p3_times_2, label="Parallel Merge (p=3)", marker='o')
plt.plot(sizes_2, p6_times_2, label="Parallel Merge (p=6)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Sequential vs Segment Merge Sort (Different Thread Levels)')
plt.legend()
plt.grid(True)
plt.savefig('plot_sequential_vs_parallel_sort.png')

# Data from the first image (Sequential vs. Segment Merge)
sizes_1 = [1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000, 512000, 1024000]
sequential_times = [3022.92, 6267.33, 13635.68, 32816.81, 59413.96, 118755.16, 237050.44, 467344.82, 937771.39, 1915051.71, 3832682.57]
segment_times_1 = [28331.16, 32812.87, 42171.91, 65323.66, 91376.39, 160360.79, 306391.84, 375526.48, 1124228.72, 1753683.13, 2747065.85]

# Data from the second image (Basic vs. Segment Merge Sort)
sizes_2 = [1000, 2000, 4000, 8000, 16000, 32000, 64000, 128000, 256000]
basic_times = [406399.19, 832904.06, 1351517.48, 2582745.27, 5063213.61, 9116846.91, 18256834.46, 35914352.82, 73582521.07]
segment_times_2 = [1195737.56, 2842908.03, 5902214.85, 9700880.14, 16649889.29, 31773088.59, 62194458.65, 124591088.80, 255035015.72]

# Plotting the first comparison
plt.figure(figsize=(10, 6))
plt.plot(sizes_1, sequential_times, label="Sequential Merge (1000 Iterations)", marker='o')
plt.plot(sizes_1, segment_times_1, label="Segment Merge (1000 Iterations)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Sequential vs Segment Merge - 1000 Iterations')
plt.legend()
plt.grid(True)
plt.savefig('plot_sequential_vs_segment_marge.png')

# Plotting the second comparison
plt.figure(figsize=(10, 6))
plt.plot(sizes_2, basic_times, label="Basic Merge Sort (100 Iterations)", marker='o')
plt.plot(sizes_2, segment_times_2, label="Segment Merge Sort (100 Iterations)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Basic vs Segment Merge Sort - 100 Iterations')
plt.legend()
plt.grid(True)

plt.savefig('plot_basic_vs_segment_sort.png')
plt.show()
