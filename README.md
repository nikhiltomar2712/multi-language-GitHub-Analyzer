# 🚀 GitHub Repository Analyzer

A powerful, multi-language GitHub repository analysis tool that provides deep insights into repository health, community engagement, and development patterns.

**Available in:** Python | TypeScript | Go | Ruby

---

## 📋 Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Implementations](#implementations)
- [Analysis Metrics](#analysis-metrics)
- [Health Scoring](#health-scoring)
- [API Rate Limits](#api-rate-limits)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

---

## ✨ Features

✅ **Multi-Language Support** - Use Python, TypeScript, Go, or Ruby
✅ **Comprehensive Analysis** - Stars, forks, contributors, languages, and more
✅ **Health Scoring** - Intelligent repository health assessment (0-100 scale)
✅ **Real-time Insights** - Actionable recommendations for improvement
✅ **High Performance** - Concurrent API requests for faster analysis
✅ **Beautiful Output** - Rich, formatted console output with emoji indicators
✅ **GitHub API Integration** - Seamless integration with GitHub's REST API v3
✅ **Token Support** - Increased rate limits with personal access tokens

---

## 🎯 Quick Start

### Using Python (Simplest)

```bash
# Install dependencies
pip install -r requirements.txt

# Set GitHub token (optional, for higher rate limits)
export GITHUB_TOKEN="ghp_your_token_here"

# Analyze a repository
python3 githubanalyzer.py torvalds/linux
```

### Using TypeScript (Node.js)

```bash
# Install dependencies
npm install

# Analyze a repository
npx ts-node github_analyzer.ts facebook/react --token ghp_xxx
```

### Using Go (Fastest)

```bash
# Build the Go binary
go build -o github-analyzer github_analyzer.go

# Analyze a repository
./github-analyzer torvalds/linux
```

### Using Ruby (Most Elegant)

```bash
# Install dependencies
gem install bundler && bundle install

# Analyze a repository
ruby github_analyzer.rb kubernetes/kubernetes
```

---

## 📦 Installation

### Prerequisites

- GitHub account (free)
- One of: Python 3.8+, Node.js 14+, Go 1.15+, or Ruby 2.7+

### Step 1: Clone or Download

```bash
git clone https://github.com/nikhiltomar2712/project.git
cd project
```

### Step 2: Choose Your Language

#### Python Setup
```bash
pip install -r requirements.txt
# or with venv
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
```

#### TypeScript Setup
```bash
npm install
# Dependencies: axios
```

#### Go Setup
```bash
# No additional dependencies needed - uses Go standard library
go build -o github-analyzer github_analyzer.go
```

#### Ruby Setup
```bash
# Ruby uses built-in libraries, no gems required
# Tested with Ruby 2.7+
```

### Step 3: Set Up GitHub Token (Optional but Recommended)

Get a personal access token from GitHub:
1. Go to [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)
2. Click "Generate new token"
3. Select `public_repo` scope (minimum required)
4. Copy the token

Set it as an environment variable:
```bash
export GITHUB_TOKEN="ghp_your_token_here"
```

---

## 🔧 Usage

### Basic Usage

```bash
# Python
python3 githubanalyzer.py owner/repository

# TypeScript
npx ts-node github_analyzer.ts owner/repository

# Go
./github-analyzer owner/repository

# Ruby
ruby github_analyzer.rb owner/repository
```

### With Token Flag

```bash
# Python
python3 githubanalyzer.py facebook/react --token ghp_xxxxx

# TypeScript
npx ts-node github_analyzer.ts facebook/react --token ghp_xxxxx

# Go
./github-analyzer facebook/react -token ghp_xxxxx

# Ruby
ruby github_analyzer.rb facebook/react --token ghp_xxxxx
```

### Example Analyses

```bash
# Analyze the Linux Kernel
python3 githubanalyzer.py torvalds/linux

# Analyze React
npx ts-node github_analyzer.ts facebook/react

# Analyze Kubernetes
./github-analyzer kubernetes/kubernetes

# Analyze Rails
ruby github_analyzer.rb rails/rails
```

---

## 🏗️ Implementations

### 📍 Python (`githubanalyzer.py` & `github_metrics.py`)

**Strengths:**
- Easiest to understand and modify
- Rich library ecosystem
- Perfect for quick analysis
- Built-in parallel fetching capability

**Speed:** ⭐⭐⭐ (≈2-3 seconds)

**Key Features:**
- `GitHubMetricsAnalyzer` class for advanced metrics
- Health scoring with recommendations
- Language distribution analysis
- Contributor statistics

```python
from github_metrics import analyze_repository_health

result = analyze_repository_health("torvalds", "linux", token="ghp_xxx")
print(result)
```

---

### 🔷 TypeScript (`github_analyzer.ts`)

**Strengths:**
- Type-safe with full TypeScript support
- Excellent for web integration
- Modern async/await patterns
- Performance monitoring

**Speed:** ⭐⭐⭐⭐ (≈1-2 seconds)

**Key Features:**
- `GitHubAnalyzerTS` class with TypeScript interfaces
- Concurrent Promise.all() for parallel requests
- Beautiful CLI output with emojis
- Health rating system

```typescript
import { GitHubAnalyzerTS } from './github_analyzer';

const analyzer = new GitHubAnalyzerTS(token);
const result = await analyzer.analyze("facebook", "react");
```

---

### 🔴 Go (`github_analyzer.go`)

**Strengths:**
- **Fastest execution** (~500ms)
- Excellent concurrency with goroutines
- Single binary deployment
- No runtime dependencies

**Speed:** ⭐⭐⭐⭐⭐ (≈0.5-1 second)

**Key Features:**
- Concurrent goroutines for parallel API calls
- `GitHubClient` with error handling
- Efficient memory usage
- Rich output formatting

```go
client := NewGitHubClient(token)
result, err := client.Analyze("kubernetes", "kubernetes")
```

---

### 💎 Ruby (`github_analyzer.rb`)

**Strengths:**
- Most elegant and readable code
- Object-oriented design with value objects
- Rails-friendly patterns
- Best for learning and teaching

**Speed:** ⭐⭐⭐ (≈2-3 seconds)

**Key Features:**
- `GitHubAnalyzer` class
- Value objects: `RepositoryInfo`, `Contributor`, `HealthScore`
- CLI interface with OptionParser
- Beautiful report formatting

```ruby
analyzer = GitHubAnalyzer.new(token)
result = analyzer.analyze("rails", "rails")
result.print_report
```

---

## 📊 Analysis Metrics

Each analysis provides:

### Repository Statistics
- **Stars** - Repository popularity
- **Forks** - Code reuse and adoption
- **Watchers** - Subscriber count
- **Open Issues** - Maintenance backlog
- **Languages** - Technology stack
- **Topics** - Project categories

### Community Metrics
- **Contributors** - Team size and diversity
- **Top Contributors** - Key project members
- **Fork-to-Star Ratio** - Community contribution level

### Activity Metrics
- **Last Update** - Maintenance recency
- **Creation Date** - Project age
- **Activity Frequency** - Maintenance status

### Derived Insights
- **Language Distribution** - Primary technologies
- **Community Health** - Engagement level
- **Maintenance Status** - Active/Stale/Archived

---

## 🏥 Health Scoring

The health score (0-100) evaluates:

| Metric | Weight | Excellent | Good | Fair | Needs Help |
|--------|--------|-----------|------|------|------------|
| Community (Stars) | 25% | 1000+ | 500+ | 100+ | < 100 |
| Contributors | 25% | 50+ | 20+ | 5+ | < 5 |
| Code Reuse (Forks) | 20% | 30%+ ratio | 10%+ ratio | - | Low |
| Issue Management | 15% | < 20 issues | < 50 issues | - | 50+ |
| Activity | 15% | < 7 days | < 30 days | < 90 days | 90+ days |

### Health Ratings

- **🟢 Excellent** (90-100): Thriving, well-maintained project
- **🟡 Good** (70-89): Healthy, active development
- **🟠 Fair** (50-69): Stable but may need attention
- **🔴 Needs Improvement** (< 50): Requires improvements

---

## 🔐 API Rate Limits

GitHub API rate limits depend on authentication:

| Type | Limit | Time Window |
|------|-------|-------------|
| **Without Token** | 60 requests | 1 hour |
| **With Token** | 5,000 requests | 1 hour |
| **Enterprise** | 15,000 requests | 1 hour |

**Tip:** Always use a token for production use.

---

## 📚 Examples

### Example 1: Analyze Linux Kernel

```bash
# Python
python3 githubanalyzer.py torvalds/linux

# Output preview:
# 📦 Repository: linux
# ⭐ Stars: 167000
# 👥 Contributors: 1200+
# 🏥 Health Score: 92/100
# Rating: Excellent
```

### Example 2: Compare Multiple Repositories

```python
# Python script to compare repos
from github_metrics import analyze_repository_health

repos = [
    "torvalds/linux",
    "rails/rails",
    "facebook/react"
]

for repo in repos:
    owner, name = repo.split("/")
    result = analyze_repository_health(owner, name, token=os.getenv("GITHUB_TOKEN"))
    print(f"{repo}: {result['health_score']}/100 - {result['contributors']} contributors")
```

### Example 3: Continuous Monitoring

```typescript
// TypeScript script to monitor repository daily
import { GitHubAnalyzerTS } from './github_analyzer';
import * as fs from 'fs';

async function monitorRepository(owner: string, repo: string) {
  const analyzer = new GitHubAnalyzerTS(process.env.GITHUB_TOKEN);
  const result = await analyzer.analyze(owner, repo);
  
  const report = {
    timestamp: new Date().toISOString(),
    health: result.health.score,
    stars: result.repository.stars,
    contributors: result.contributors.length
  };
  
  fs.appendFileSync('monitoring.json', JSON.stringify(report) + '\n');
}

// Run daily with cron
// 0 0 * * * /usr/bin/node monitor.js
```

---

## 🛠️ Development

### Building from Source

```bash
# Python - No build needed, just install dependencies
pip install -r requirements.txt

# TypeScript - Compile to JavaScript
npm run build

# Go - Build optimized binary
go build -ldflags="-s -w" -o github-analyzer github_analyzer.go

# Ruby - No build needed, runs directly
ruby github_analyzer.rb --help
```

### Running Tests

```bash
# Python
pytest tests/

# TypeScript
npm test

# Go
go test ./...

# Ruby
rspec spec/
```

---

## 🐛 Troubleshooting

### "Repository not found" Error

**Cause:** Invalid repository name or no internet connection
**Solution:** Check repository name format and internet connection

```bash
# Correct format:
python3 githubanalyzer.py facebook/react

# Wrong format:
python3 githubanalyzer.py facebook-react  # ❌ Incorrect
```

### "Rate limit exceeded" Error

**Cause:** Too many requests without authentication
**Solution:** Use GitHub token to increase limit from 60 to 5,000 requests/hour

```bash
export GITHUB_TOKEN="ghp_your_token_here"
python3 githubanalyzer.py torvalds/linux  # ✅ Now works
```

### "Unauthorized" Error (401)

**Cause:** Invalid or expired GitHub token
**Solution:** Generate new token from [GitHub settings](https://github.com/settings/tokens)

```bash
# Verify token
curl -H "Authorization: token YOUR_TOKEN" https://api.github.com/user
```

### Connection Timeout

**Cause:** Slow network or GitHub API unavailable
**Solution:** Increase timeout or retry after few seconds

```python
# Python - modify timeout in github_metrics.py
requests.Timeout(60)
```

---

## 📈 Performance Comparison

```
Repository: torvalds/linux

Tool         | Time  | Language | Dependencies
-------------|-------|----------|------------------
Go binary    | ~0.5s | Go       | None (Stdlib)
TypeScript   | ~1.2s | Node.js  | axios
Python       | ~2.1s | Python   | requests
Ruby         | ~2.3s | Ruby     | Stdlib
```

**Recommendation:**
- **Speed critical:** Use Go
- **Web apps:** Use TypeScript
- **Scripting:** Use Python or Ruby
- **Learning:** Use Ruby

---

## 🤝 Contributing

Contributions are welcome! Areas for improvement:

1. **Additional Metrics**
   - Code quality metrics (via Code Climate API)
   - Test coverage
   - Dependency analysis

2. **New Languages**
   - Rust implementation
   - Java implementation
   - C# implementation

3. **Features**
   - Batch analysis
   - JSON export
   - Database storage
   - Web dashboard
   - GitHub Action integration

4. **Documentation**
   - API documentation
   - More examples
   - Video tutorials

### Steps to Contribute

```bash
# 1. Fork the repository
# 2. Create a feature branch
git checkout -b feature/new-feature

# 3. Make changes and commit
git commit -m "Add new feature"

# 4. Push to branch
git push origin feature/new-feature

# 5. Create Pull Request
```

---

## 📄 License

MIT License - See LICENSE file for details

This project is free and open-source. You can use it for personal and commercial projects.

---

## 🔗 Resources

- [GitHub API Documentation](https://docs.github.com/rest)
- [Personal Access Tokens](https://github.com/settings/tokens)
- [GitHub API Rate Limiting](https://docs.github.com/rest/overview/resources-in-the-rest-api#rate-limiting)
- [Repository Insights Guide](https://docs.github.com/en/repositories/viewing-activity-and-data-for-your-repository)

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/nikhiltomar2712/project/issues)
- **Discussions:** [GitHub Discussions](https://github.com/nikhiltomar2712/project/discussions)
- **Email:** [Contact](mailto:nikhiltomarsan2712@gmail.com)

---

## 🌟 Star History

Give this project a ⭐ if you find it useful!

---

## 👨‍💻 Author

Created by **Nikhil Tomar**
- GitHub: [@nikhiltomar2712](https://github.com/nikhiltomar2712)
- Portfolio: [nikhiltomar.dev](https://nikhiltomar.dev)

---

**Last Updated:** June 2026
**Version:** 2.1.0
**Status:** ✅ Actively Maintained
