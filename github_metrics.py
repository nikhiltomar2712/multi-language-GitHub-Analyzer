#!/usr/bin/env python3
"""
Enhanced GitHub Metrics Module
Provides advanced analytics for GitHub repositories including contributor insights,
language distribution, and activity metrics.
"""

import requests
from datetime import datetime
from typing import Dict, List, Optional
from collections import defaultdict

class GitHubMetricsAnalyzer:
    """Advanced metrics analyzer for GitHub repositories"""
    
    def __init__(self, token: Optional[str] = None):
        self.token = token
        self.base_url = "https://api.github.com"
        self.headers = self._build_headers()
    
    def _build_headers(self) -> Dict[str, str]:
        """Build request headers with authentication"""
        headers = {"Accept": "application/vnd.github.v3+json"}
        if self.token:
            headers["Authorization"] = f"token {self.token}"
        return headers
    
    def get_language_distribution(self, owner: str, repo: str) -> Dict[str, int]:
        """
        Get programming language distribution in repository
        
        Args:
            owner: Repository owner
            repo: Repository name
            
        Returns:
            Dictionary of languages and their byte counts
        """
        url = f"{self.base_url}/repos/{owner}/{repo}/languages"
        response = requests.get(url, headers=self.headers)
        response.raise_for_status()
        return response.json()
    
    def get_contributor_stats(self, owner: str, repo: str) -> List[Dict]:
        """
        Get detailed contributor statistics
        
        Args:
            owner: Repository owner
            repo: Repository name
            
        Returns:
            List of contributor data with commit counts
        """
        url = f"{self.base_url}/repos/{owner}/{repo}/contributors"
        response = requests.get(url, headers=self.headers, params={"per_page": 100})
        response.raise_for_status()
        return response.json()
    
    def get_repository_info(self, owner: str, repo: str) -> Dict:
        """
        Get comprehensive repository information
        
        Args:
            owner: Repository owner
            repo: Repository name
            
        Returns:
            Repository metadata and statistics
        """
        url = f"{self.base_url}/repos/{owner}/{repo}"
        response = requests.get(url, headers=self.headers)
        response.raise_for_status()
        data = response.json()
        
        return {
            "name": data.get("name"),
            "description": data.get("description"),
            "url": data.get("html_url"),
            "stars": data.get("stargazers_count"),
            "forks": data.get("forks_count"),
            "open_issues": data.get("open_issues_count"),
            "watchers": data.get("watchers_count"),
            "created_at": data.get("created_at"),
            "updated_at": data.get("updated_at"),
            "language": data.get("language"),
            "topics": data.get("topics"),
            "is_fork": data.get("fork"),
            "is_archived": data.get("archived")
        }
    
    def get_activity_summary(self, owner: str, repo: str) -> Dict:
        """
        Get repository activity summary
        
        Args:
            owner: Repository owner
            repo: Repository name
            
        Returns:
            Activity metrics and trends
        """
        url = f"{self.base_url}/repos/{owner}/{repo}/commits"
        response = requests.get(url, headers=self.headers, params={"per_page": 1})
        response.raise_for_status()
        
        latest_commit = response.json()[0] if response.json() else None
        
        return {
            "latest_commit": latest_commit.get("commit").get("author").get("date") if latest_commit else None,
            "latest_commit_sha": latest_commit.get("sha") if latest_commit else None,
            "latest_commit_author": latest_commit.get("commit").get("author").get("name") if latest_commit else None
        }


def analyze_repository_health(owner: str, repo: str, token: Optional[str] = None) -> Dict:
    """
    Analyze overall repository health
    
    Args:
        owner: Repository owner
        repo: Repository name
        token: GitHub API token (optional)
        
    Returns:
        Health score and recommendations
    """
    analyzer = GitHubMetricsAnalyzer(token)
    
    try:
        info = analyzer.get_repository_info(owner, repo)
        contributors = analyzer.get_contributor_stats(owner, repo)
        languages = analyzer.get_language_distribution(owner, repo)
        activity = analyzer.get_activity_summary(owner, repo)
        
        # Calculate health score
        health_score = 0
        recommendations = []
        
        if info.get("stars", 0) > 100:
            health_score += 25
        else:
            recommendations.append("Increase community engagement for more stars")
        
        if len(contributors) > 5:
            health_score += 25
        else:
            recommendations.append("Encourage more contributors")
        
        if len(languages) > 1:
            health_score += 20
        else:
            recommendations.append("Consider adding supporting languages")
        
        if info.get("open_issues", 0) < 50:
            health_score += 20
        else:
            recommendations.append("Review and close stale issues")
        
        if not info.get("is_archived"):
            health_score += 10
        else:
            recommendations.append("Consider archiving inactive repositories")
        
        return {
            "health_score": health_score,
            "repository": info.get("name"),
            "contributors": len(contributors),
            "languages": list(languages.keys()),
            "stars": info.get("stars"),
            "recommendations": recommendations,
            "last_activity": activity.get("latest_commit")
        }
    except requests.exceptions.RequestException as e:
        return {"error": str(e)}


if __name__ == "__main__":
    import sys
    import os
    
    if len(sys.argv) < 2:
        print("Usage: python github_metrics.py <owner/repo> [--token TOKEN]")
        sys.exit(1)
    
    repo_path = sys.argv[1]
    token = os.getenv("GITHUB_TOKEN")
    
    if "--token" in sys.argv:
        token_idx = sys.argv.index("--token")
        if token_idx + 1 < len(sys.argv):
            token = sys.argv[token_idx + 1]
    
    owner, repo = repo_path.split("/")
    result = analyze_repository_health(owner, repo, token)
    
    import json
    print(json.dumps(result, indent=2))
