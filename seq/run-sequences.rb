#!/usr/bin/ruby

DIR = File.dirname(__FILE__)

def main
  n_begin = ARGV.shift.to_i
  n_end = ARGV.shift.to_i
  sequence_names = ARGV

  data_filename = data_filename_for_sequences(sequence_names)

  log_filename = log_filename_for_sequences(sequence_names)

  puts "writing #{data_filename} and #{log_filename}"

  run_sequence(n_begin, n_end, sequence_names, data_filename, log_filename)
end

def data_filename_for_sequences(sequence_names)
  sequence_names.join('-') + '.txt'
end

def log_filename_for_sequences(sequence_names)
  sequence_names.join('-') + '.log'
end

def run_sequence(n_begin, n_end, sequence_names, data_filename, log_filename)
  sys("#{DIR}/../bin/sequence -b #{n_begin} -e #{n_end} #{sequence_names.join(' ')} > #{data_filename} 2> #{log_filename}")
end

def sys(cmd)
  puts "exec: #{cmd}"
  system(cmd) || raise
end

main
