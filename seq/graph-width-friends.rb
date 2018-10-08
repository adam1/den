#!/usr/bin/ruby

def make_plot(xrange, logscale_y)
  sequences_spec = {
    "WidthV3" => "#ff0000",
    "NumTypes" => "#0000ff",
  }

  graph_filename = graph_filename_for_sequences(sequences_spec, xrange, logscale_y)

  puts "writing #{graph_filename}"

  plot_clauses = plot_clauses_for_sequences(sequences_spec, xrange)

  xrange_clause = xrange ? "set xrange [#{xrange[0]}:#{xrange[1]}]" : ""

  yrange_clause = ""
  logscale_y_clause = ""

  if logscale_y
    yrange_clause = xrange && xrange[0] == 70 ? "set yrange [2e+90:2e+120]" : ""
    logscale_y_clause = "set logscale y 2"
  end

  plot_cmds = <<END
set terminal png
set output "#{graph_filename}"
set grid xtics lt 0
set grid ytics lt 0
#{logscale_y_clause}
set key center right
#{xrange_clause}
#{yrange_clause}
#unset logscale y
#set format y "%.0s%cB"
f(x) = int(x)!
plot #{plot_clauses}
END

  #puts plot_cmds
  run_gnuplot(plot_cmds)
end

def plot_clauses_for_sequences(sequences_spec, xrange)
  plot_clauses = sequences_spec.map {|name, color| plot_clause_for_sequence(name, color, xrange)}
  #plot_clauses << %Q(f(x) with lines linecolor rgbcolor "#00ff00" title "n!")
  #plot_clauses << "g(x)"
  plot_clauses.join(', ')
end

def plot_clause_for_sequence(name, color, xrange)
  data_filename = name + '.txt'
  style = xrange && xrange[0] > 0 ? "linespoints pt 5" : "lines"
  %Q("#{data_filename}" using 1:2 with #{style} linecolor rgbcolor "#{color}" title "#{name}")
end

def graph_filename_for_sequences(sequences_spec, xrange, logscale_y)
  fn = sequences_spec.keys.sort.join('-')
  if xrange && xrange[0] > 0
    fn += '-' + xrange.join('-')
  end
  if logscale_y
    fn += '-log'
  end
  fn + '.png'
end

def run_gnuplot(cmds)
  system("echo '#{cmds}' | gnuplot") || raise
end

def main
#   make_plot([0, 80], true)
#   make_plot([40, 80], true)
#   make_plot([70, 80], true)
#   make_plot([78, 80], true)

  make_plot([0, 80], false)
  make_plot([40, 80], false)
  make_plot([78, 80], false)
end

main
