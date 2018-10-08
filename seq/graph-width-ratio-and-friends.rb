#!/usr/bin/ruby

def main
  sequence_names = ["WidthV3SuccessiveRatio"]

  graph_filename = graph_filename_for_sequences(sequence_names)

  puts "writing #{graph_filename}"

  plot_clauses = plot_clauses_for_sequences(sequence_names)

  plot_cmds = <<END
set terminal png
set output "#{graph_filename}"
set grid xtics lt 0
set grid ytics lt 0
#set logscale y 2
#unset logscale y
#set format y "%.0s%cB"
f(x) = x
#g(x) = ((int(x - 1))!)/x
plot #{plot_clauses}
END

  #puts plot_cmds
  run_gnuplot(plot_cmds)
end

def plot_clauses_for_sequences(sequence_names)
  plot_clauses = sequence_names.map {|seq| plot_clause_for_sequence(seq)}
  plot_clauses << %Q(f(x) title "n")
  #plot_clauses << "g(x)"
  plot_clauses.join(', ')
end

def plot_clause_for_sequence(sequence_name)
  data_filename = sequence_name + '.txt'
  %Q("#{data_filename}" using 1:2 with lines title "#{sequence_name}")
end

def graph_filename_for_sequences(sequence_names)
  sequence_names.join('-') + '.png'
end

def run_gnuplot(cmds)
  system("echo '#{cmds}' | gnuplot") || raise
end

main
