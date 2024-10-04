import matplotlib.pyplot as plt
import pandas as pd
import numpy as np


merge_file_path = './merge_data.csv'
merge_data = pd.read_csv(merge_file_path)

sort_file_path = './sort_data.csv'
sort_data = pd.read_csv(sort_file_path)



# Merge data
merge_sizes = merge_data['SIZE']
merge_sequential_times = merge_data['Sequential']
merge_p1_times = merge_data['p=1']
merge_p3_times = merge_data['p=3']
merge_p6_times = merge_data['p=6']

# Sort data
sort_sizes = sort_data['SIZE']
sort_basic_times = sort_data['Basic']
sort_p1_times = sort_data['p=1']
sort_p3_times = sort_data['p=3']
sort_p6_times = sort_data['p=6']

# Plotting the first comparison (Sequential vs Segment)
plt.figure(figsize=(10, 6))
plt.plot(merge_sizes, merge_sequential_times, label="Sequential Merge", marker='o')
plt.plot(merge_sizes, merge_p1_times, label="Segment Merge (p=1)", marker='o')
plt.plot(merge_sizes, merge_p3_times, label="Segment Merge (p=3)", marker='o')
plt.plot(merge_sizes, merge_p6_times, label="Segment Merge (p=6)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Sequential vs Segment Merge (100 Iteration Average)')
plt.legend()
plt.grid(True)
plt.savefig('plot_merges.png')

# Plotting the second comparison (Basic vs Parallel)
plt.figure(figsize=(10, 6))
plt.plot(sort_sizes, sort_basic_times, label="Basic Merge Sort", marker='o')
plt.plot(sort_sizes, sort_p1_times, label="Fully Parallel Merge Sort (p=1)", marker='o')
plt.plot(sort_sizes, sort_p3_times, label="Fully Parallel Merge Sort (p=3)", marker='o')
plt.plot(sort_sizes, sort_p6_times, label="Fully Parallel Merge Sort (p=6)", marker='o')
plt.xscale('log')
plt.yscale('log')
plt.xlabel('Input Size (log scale)')
plt.ylabel('Time (ns) (log scale)')
plt.title('Basic vs Fully Parallel Merge Sort (100 Iteration Average)')
plt.legend()
plt.grid(True)
plt.savefig('plot_sorts.png')
