#!/usr/bin/ruby

def main
  x_start = ARGV[0].to_i
  graph_filename = "Density-friends-#{x_start}.png"

  puts "writing #{graph_filename}"

  plot_cmds = <<END
set terminal png
set output "#{graph_filename}"
set grid xtics lt 0
set grid ytics lt 0
set xrange [#{x_start}:]
f1(x) = 1/x
plot f1(x) with lines title "1/x", \
  "Density.txt" using 1:2 with lines title "Density"
END

  #puts plot_cmds
  run_gnuplot(plot_cmds)
  system("open #{graph_filename}")
end

def run_gnuplot(cmds)
  system("echo '#{cmds}' | gnuplot") || raise
end

main
