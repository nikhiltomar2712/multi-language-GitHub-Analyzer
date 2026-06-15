#!/usr/bin/env node
/**
 * GitHub Analyzer CLI - TypeScript Version
 * A modern TypeScript implementation for analyzing GitHub repositories
 * with type safety and concurrent request handling
 */

import axios, { AxiosInstance } from 'axios';
import { performance } from 'perf_hooks';

interface GitHubRepo {
  name: string;
  description: string | null;
  stars: number;
  forks: number;
  watchers: number;
  language: string | null;
  topics: string[];
  createdAt: string;
  updatedAt: string;
  openIssues: number;
}

interface Contributor {
  login: string;
  contributions: number;
  avatarUrl: string;
}

interface AnalysisResult {
  repository: GitHubRepo;
  contributors: Contributor[];
  languages: { [key: string]: number };
  health: RepositoryHealth;
  analysisTime: number;
}

interface RepositoryHealth {
  score: number;
  rating: 'Excellent' | 'Good' | 'Fair' | 'Needs Improvement';
  insights: string[];
}

class GitHubAnalyzerTS {
  private client: AxiosInstance;
  private token: string | null;

  constructor(token?: string) {
    this.token = token || process.env.GITHUB_TOKEN || null;
    this.client = axios.create({
      baseURL: 'https://api.github.com',
      headers: {
        'Accept': 'application/vnd.github.v3+json',
        ...(this.token && { 'Authorization': `token ${this.token}` })
      }
    });
  }

  async getRepositoryInfo(owner: string, repo: string): Promise<GitHubRepo> {
    try {
      const response = await this.client.get(`/repos/${owner}/${repo}`);
      const data = response.data;

      return {
        name: data.name,
        description: data.description,
        stars: data.stargazers_count,
        forks: data.forks_count,
        watchers: data.watchers_count,
        language: data.language,
        topics: data.topics || [],
        createdAt: data.created_at,
        updatedAt: data.updated_at,
        openIssues: data.open_issues_count
      };
    } catch (error) {
      throw new Error(`Failed to fetch repository info: ${error}`);
    }
  }

  async getContributors(owner: string, repo: string): Promise<Contributor[]> {
    try {
      const response = await this.client.get(`/repos/${owner}/${repo}/contributors`, {
        params: { per_page: 50 }
      });

      return response.data.map((contributor: any) => ({
        login: contributor.login,
        contributions: contributor.contributions,
        avatarUrl: contributor.avatar_url
      }));
    } catch (error) {
      throw new Error(`Failed to fetch contributors: ${error}`);
    }
  }

  async getLanguages(owner: string, repo: string): Promise<{ [key: string]: number }> {
    try {
      const response = await this.client.get(`/repos/${owner}/${repo}/languages`);
      return response.data;
    } catch (error) {
      throw new Error(`Failed to fetch languages: ${error}`);
    }
  }

  calculateHealth(repo: GitHubRepo, contributors: Contributor[]): RepositoryHealth {
    let score = 0;
    const insights: string[] = [];

    // Star metric (0-25 points)
    if (repo.stars >= 1000) {
      score += 25;
      insights.push('🌟 Excellent community engagement');
    } else if (repo.stars >= 500) {
      score += 20;
      insights.push('⭐ Good community interest');
    } else if (repo.stars >= 100) {
      score += 15;
      insights.push('📈 Growing community');
    } else {
      insights.push('🚀 Building community awareness');
    }

    // Contributor metric (0-25 points)
    if (contributors.length >= 50) {
      score += 25;
      insights.push('👥 Strong contributor base');
    } else if (contributors.length >= 20) {
      score += 20;
      insights.push('🤝 Active community');
    } else if (contributors.length >= 5) {
      score += 15;
      insights.push('👤 Growing team');
    } else {
      insights.push('👤 Building contributor network');
    }

    // Fork metric (0-20 points)
    if (repo.forks >= repo.stars * 0.3) {
      score += 20;
      insights.push('🔀 High fork-to-star ratio');
    } else if (repo.forks >= repo.stars * 0.1) {
      score += 15;
      insights.push('🔀 Moderate fork activity');
    } else {
      score += 10;
    }

    // Issue metric (0-15 points)
    if (repo.openIssues < 20) {
      score += 15;
      insights.push('✅ Well-maintained issue tracker');
    } else if (repo.openIssues < 50) {
      score += 10;
      insights.push('📋 Manageable issue backlog');
    } else {
      insights.push('⚠️ Consider addressing issue backlog');
    }

    // Activity metric (0-15 points)
    const lastUpdate = new Date(repo.updatedAt).getTime();
    const daysSinceUpdate = (Date.now() - lastUpdate) / (1000 * 60 * 60 * 24);
    
    if (daysSinceUpdate < 7) {
      score += 15;
      insights.push('🔥 Actively maintained');
    } else if (daysSinceUpdate < 30) {
      score += 10;
      insights.push('✓ Regularly maintained');
    } else if (daysSinceUpdate < 90) {
      score += 5;
      insights.push('⏱️ Maintenance may be slowing');
    }

    const rating = this.scoreToRating(score);

    return { score, rating, insights };
  }

  private scoreToRating(score: number): 'Excellent' | 'Good' | 'Fair' | 'Needs Improvement' {
    if (score >= 90) return 'Excellent';
    if (score >= 70) return 'Good';
    if (score >= 50) return 'Fair';
    return 'Needs Improvement';
  }

  async analyze(owner: string, repo: string): Promise<AnalysisResult> {
    const startTime = performance.now();

    const [repoInfo, contributors, languages] = await Promise.all([
      this.getRepositoryInfo(owner, repo),
      this.getContributors(owner, repo),
      this.getLanguages(owner, repo)
    ]);

    const health = this.calculateHealth(repoInfo, contributors);
    const analysisTime = Math.round((performance.now() - startTime) * 100) / 100;

    return {
      repository: repoInfo,
      contributors,
      languages,
      health,
      analysisTime
    };
  }
}

// CLI Interface
async function main() {
  const args = process.argv.slice(2);

  if (args.length === 0 || args.includes('--help') || args.includes('-h')) {
    console.log(`
GitHub Analyzer CLI - TypeScript Version
Usage: npx github-analyzer-ts <owner/repo> [options]

Options:
  --token TOKEN    GitHub API token (or set GITHUB_TOKEN environment variable)
  --help, -h       Show this help message

Example:
  npx github-analyzer-ts torvalds/linux
  npx github-analyzer-ts facebook/react --token ghp_xxxxx
    `);
    process.exit(0);
  }

  const repoPath = args[0];
  let token: string | undefined;

  if (args.includes('--token')) {
    const tokenIdx = args.indexOf('--token');
    if (tokenIdx + 1 < args.length) {
      token = args[tokenIdx + 1];
    }
  }

  try {
    const [owner, repo] = repoPath.split('/');
    if (!owner || !repo) {
      throw new Error('Invalid repository format. Use: owner/repo');
    }

    const analyzer = new GitHubAnalyzerTS(token);
    console.log(`\n📊 Analyzing ${owner}/${repo}...\n`);

    const result = await analyzer.analyze(owner, repo);

    // Display results
    console.log('╔════════════════════════════════════════════════════════╗');
    console.log('║           GITHUB REPOSITORY ANALYSIS REPORT            ║');
    console.log('╚════════════════════════════════════════════════════════╝\n');

    console.log(`📦 Repository: ${result.repository.name}`);
    if (result.repository.description) {
      console.log(`📝 Description: ${result.repository.description}`);
    }

    console.log(`\n📊 Statistics:`);
    console.log(`  ⭐ Stars: ${result.repository.stars}`);
    console.log(`  🔀 Forks: ${result.repository.forks}`);
    console.log(`  👁️ Watchers: ${result.repository.watchers}`);
    console.log(`  ⚠️ Open Issues: ${result.repository.openIssues}`);

    console.log(`\n👥 Contributors: ${result.contributors.length}`);
    const topContributors = result.contributors.slice(0, 5);
    topContributors.forEach((c, idx) => {
      console.log(`  ${idx + 1}. ${c.login} (${c.contributions} commits)`);
    });

    console.log(`\n💻 Languages:`);
    Object.entries(result.languages)
      .sort(([, a], [, b]) => b - a)
      .slice(0, 5)
      .forEach(([lang, bytes]) => {
        const percent = ((bytes / Object.values(result.languages).reduce((a, b) => a + b, 0)) * 100).toFixed(1);
        console.log(`  ${lang}: ${percent}%`);
      });

    console.log(`\n🏥 Health Assessment:`);
    console.log(`  Health Score: ${result.health.score}/100`);
    console.log(`  Rating: ${result.health.rating}`);
    console.log(`\n  Insights:`);
    result.health.insights.forEach(insight => {
      console.log(`  ${insight}`);
    });

    console.log(`\n⏱️ Analysis completed in ${result.analysisTime}ms\n`);

  } catch (error) {
    console.error(`❌ Error: ${error instanceof Error ? error.message : String(error)}`);
    process.exit(1);
  }
}

main();

export { GitHubAnalyzerTS, GitHubRepo, Contributor, AnalysisResult };
