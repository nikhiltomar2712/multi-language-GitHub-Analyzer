#!/usr/bin/env ruby
# frozen_string_literal: true

require 'net/http'
require 'json'
require 'optparse'
require 'time'

##
# GitHub Analyzer - Ruby Implementation
# A elegant and idiomatic Ruby approach to analyzing GitHub repositories
#

class GitHubAnalyzer
  BASE_URL = 'https://api.github.com'
  DEFAULT_TIMEOUT = 30
  
  attr_reader :token, :http_client

  def initialize(token = nil)
    @token = token || ENV['GITHUB_TOKEN']
    @http_client = Net::HTTP
  end

  ##
  # Fetch repository information from GitHub API
  #
  def get_repository(owner, repo)
    endpoint = "/repos/#{owner}/#{repo}"
    response = make_request(endpoint)
    
    RepositoryInfo.new(
      name: response['name'],
      description: response['description'],
      stars: response['stargazers_count'],
      forks: response['forks_count'],
      watchers: response['watchers_count'],
      language: response['language'],
      topics: response['topics'] || [],
      created_at: response['created_at'],
      updated_at: response['updated_at'],
      open_issues: response['open_issues_count'],
      url: response['html_url'],
      is_archived: response['archived'],
      is_fork: response['fork']
    )
  end

  ##
  # Fetch repository contributors
  #
  def get_contributors(owner, repo)
    endpoint = "/repos/#{owner}/#{repo}/contributors?per_page=100"
    response = make_request(endpoint)
    
    response.map do |contributor|
      Contributor.new(
        login: contributor['login'],
        contributions: contributor['contributions'],
        avatar_url: contributor['avatar_url']
      )
    end
  end

  ##
  # Fetch programming language distribution
  #
  def get_languages(owner, repo)
    endpoint = "/repos/#{owner}/#{repo}/languages"
    response = make_request(endpoint)
    response || {}
  end

  ##
  # Perform complete analysis of a repository
  #
  def analyze(owner, repo)
    start_time = Time.now
    
    # Fetch all data concurrently (simulated with sequential calls)
    repository = get_repository(owner, repo)
    contributors = get_contributors(owner, repo)
    languages = get_languages(owner, repo)
    
    health_score = calculate_health(repository, contributors)
    duration = Time.now - start_time
    
    AnalysisResult.new(
      repository: repository,
      contributors: contributors,
      languages: languages,
      health: health_score,
      duration: duration
    )
  rescue StandardError => e
    raise "Analysis failed: #{e.message}"
  end

  private

  ##
  # Make HTTP request to GitHub API
  #
  def make_request(endpoint)
    uri = URI("#{BASE_URL}#{endpoint}")
    
    http = @http_client.new(uri.host, uri.port)
    http.use_ssl = true
    http.read_timeout = DEFAULT_TIMEOUT
    
    request = Net::HTTP::Get.new(uri.request_uri)
    request['Accept'] = 'application/vnd.github.v3+json'
    request['Authorization'] = "token #{@token}" if @token
    
    response = http.request(request)
    
    case response.code.to_i
    when 200
      JSON.parse(response.body)
    when 404
      raise "Repository not found"
    when 401
      raise "Unauthorized - check your GitHub token"
    when 403
      raise "Rate limit exceeded"
    else
      raise "HTTP #{response.code}: #{response.body}"
    end
  end

  ##
  # Calculate health score for repository
  #
  def calculate_health(repository, contributors)
    score = 0
    insights = []

    # Stars evaluation (0-25)
    case repository.stars
    when 1000..Float::INFINITY
      score += 25
      insights << "🌟 Excellent community engagement"
    when 500...1000
      score += 20
      insights << "⭐ Good community interest"
    when 100...500
      score += 15
      insights << "📈 Growing community"
    else
      insights << "🚀 Building community awareness"
    end

    # Contributors evaluation (0-25)
    case contributors.length
    when 50..Float::INFINITY
      score += 25
      insights << "👥 Strong contributor base"
    when 20...50
      score += 20
      insights << "🤝 Active community"
    when 5...20
      score += 15
      insights << "👤 Growing team"
    else
      insights << "👤 Building contributor network"
    end

    # Forks evaluation (0-20)
    if repository.stars > 0
      fork_ratio = repository.forks.to_f / repository.stars
      case fork_ratio
      when 0.3..Float::INFINITY
        score += 20
        insights << "🔀 High fork-to-star ratio"
      when 0.1...0.3
        score += 15
        insights << "🔀 Moderate fork activity"
      else
        score += 10
      end
    end

    # Issues evaluation (0-15)
    case repository.open_issues
    when 0...20
      score += 15
      insights << "✅ Well-maintained issue tracker"
    when 20...50
      score += 10
      insights << "📋 Manageable issue backlog"
    else
      insights << "⚠️ Consider addressing issue backlog"
    end

    # Activity evaluation (0-15)
    updated_at = Time.parse(repository.updated_at)
    days_since_update = (Time.now - updated_at) / 86400.0

    case days_since_update
    when 0...7
      score += 15
      insights << "🔥 Actively maintained"
    when 7...30
      score += 10
      insights << "✓ Regularly maintained"
    when 30...90
      score += 5
      insights << "⏱️ Maintenance may be slowing"
    end

    rating = case score
             when 90..100 then "Excellent"
             when 70...90 then "Good"
             when 50...70 then "Fair"
             else "Needs Improvement"
             end

    HealthScore.new(score: score, rating: rating, insights: insights)
  end
end

##
# Value object for repository information
#
class RepositoryInfo
  attr_reader :name, :description, :stars, :forks, :watchers,
              :language, :topics, :created_at, :updated_at,
              :open_issues, :url, :is_archived, :is_fork

  def initialize(**attrs)
    attrs.each { |key, value| instance_variable_set("@#{key}", value) }
  end
end

##
# Value object for contributor information
#
class Contributor
  attr_reader :login, :contributions, :avatar_url

  def initialize(login:, contributions:, avatar_url:)
    @login = login
    @contributions = contributions
    @avatar_url = avatar_url
  end
end

##
# Value object for health score
#
class HealthScore
  attr_reader :score, :rating, :insights

  def initialize(score:, rating:, insights:)
    @score = score
    @rating = rating
    @insights = insights
  end
end

##
# Value object for analysis result
#
class AnalysisResult
  attr_reader :repository, :contributors, :languages, :health, :duration

  def initialize(repository:, contributors:, languages:, health:, duration:)
    @repository = repository
    @contributors = contributors
    @languages = languages
    @health = health
    @duration = duration
  end

  ##
  # Pretty print analysis results
  #
  def print_report
    puts "\n╔════════════════════════════════════════════════════════╗"
    puts "║           GITHUB REPOSITORY ANALYSIS REPORT            ║"
    puts "╚════════════════════════════════════════════════════════╝\n"

    puts "📦 Repository: #{repository.name}"
    puts "📝 Description: #{repository.description}" if repository.description

    puts "\n📊 Statistics:"
    puts "  ⭐ Stars: #{repository.stars}"
    puts "  🔀 Forks: #{repository.forks}"
    puts "  👁️ Watchers: #{repository.watchers}"
    puts "  ⚠️ Open Issues: #{repository.open_issues}"

    puts "\n👥 Top Contributors: #{contributors.length} total"
    contributors.take(5).each_with_index do |contributor, index|
      puts "  #{index + 1}. #{contributor.login} (#{contributor.contributions} commits)"
    end

    print_languages

    print_health

    puts "\n⏱️ Analysis completed in #{(duration * 1000).round(2)}ms\n"
  end

  private

  def print_languages
    return if languages.empty?

    puts "\n💻 Languages:"
    total_bytes = languages.values.sum
    languages.sort_by { |_, bytes| -bytes }.take(5).each do |lang, bytes|
      percent = ((bytes.to_f / total_bytes) * 100).round(1)
      puts "  #{lang}: #{percent}%"
    end
  end

  def print_health
    puts "\n🏥 Health Assessment:"
    puts "  Health Score: #{health.score}/100"
    puts "  Rating: #{health.rating}"
    puts "\n  Insights:"
    health.insights.each { |insight| puts "  #{insight}" }
  end
end

##
# CLI interface
#
class CLI
  def self.run
    options = {}

    OptionParser.new do |opts|
      opts.banner = "GitHub Analyzer - Ruby Version\nUsage: ruby github_analyzer.rb <owner/repo> [options]"

      opts.on('-t', '--token TOKEN', 'GitHub API token') do |token|
        options[:token] = token
      end

      opts.on('-h', '--help', 'Show this help message') do
        puts opts
        exit
      end
    end.parse!

    if ARGV.empty?
      puts "Usage: ruby github_analyzer.rb <owner/repo> [--token TOKEN]"
      puts "Set GITHUB_TOKEN environment variable to avoid using --token flag"
      exit 1
    end

    repo_path = ARGV[0]
    owner, repo = repo_path.split('/')

    unless owner && repo
      puts "❌ Error: Invalid repository format. Use: owner/repo"
      exit 1
    end

    analyzer = GitHubAnalyzer.new(options[:token])
    puts "\n📊 Analyzing #{owner}/#{repo}...\n"

    result = analyzer.analyze(owner, repo)
    result.print_report
  rescue StandardError => e
    puts "❌ Error: #{e.message}"
    exit 1
  end
end

# Run CLI if executed directly
CLI.run if __FILE__ == $PROGRAM_NAME
