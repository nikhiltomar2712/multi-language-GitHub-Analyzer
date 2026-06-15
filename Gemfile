source 'https://rubygems.org'

git_source(:github) { |repo| "https://github.com/#{repo}.git" }

ruby '2.7.0', :engine => 'ruby', :engine_version => '2.7.0'

# No external gems required for basic functionality
# All utilities use Ruby standard library (Net::HTTP, JSON, OptionParser, Time)

# Development dependencies
group :development, :test do
  # Testing
  gem 'rspec', '~> 3.12'
  gem 'rspec-core', '~> 3.12'
  
  # Code quality
  gem 'rubocop', '~> 1.50'
  gem 'rubocop-performance', '~> 1.17'
  
  # Documentation
  gem 'yard', '~> 0.9'
  
  # Debugging
  gem 'pry', '~> 0.14'
  gem 'pry-byebug', '~> 3.10'
end

group :development do
  # Development utilities
  gem 'bundler', '~> 2.0'
  gem 'rake', '~> 13.0'
end

# Optional: For enhanced functionality
group :optional do
  # For colored output (optional)
  gem 'colorize', '~> 0.8.1'
  
  # For better error handling (optional)
  gem 'httparty', '~> 0.21.0'  # Alternative to Net::HTTP
  
  # For caching (optional)
  gem 'redis', '~> 5.0'
end

# Performance monitoring (optional)
# gem 'benchmark-ips', '~> 2.13'
