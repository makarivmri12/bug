from typing import Dict, List, Optional, Tuple
import hashlib
import difflib
import json
from dataclasses import dataclass
from enum import Enum


class VulnerabilityType(Enum):
    """Types of vulnerabilities detected"""
    IDOR = "IDOR"
    SQL_INJECTION = "SQLi"
    XSS = "XSS"
    SSRF = "SSRF"
    LOGIC_FLAW = "LogicFlaw"
    AUTH_BYPASS = "AuthBypass"
    BROKEN_ACCESS = "BrokenAccess"
    TIMING_ATTACK = "TimingAttack"
    RESPONSE_DIFF = "ResponseDiff"


@dataclass
class ResponseMetrics:
    """Metrics extracted from HTTP responses"""
    status_code: int
    content_hash: str
    content_length: int
    response_time_ms: int
    headers_hash: str
    body_preview: str


class DifferentialAnalyzer:
    """
    Advanced differential analysis engine.
    Compares baseline and mutated responses to identify anomalies.
    """

    def __init__(self, threshold_similarity: float = 0.85):
        """
        Initialize the differential analyzer.
        
        Args:
            threshold_similarity: Similarity threshold for detecting differences (0-1)
        """
        self.threshold_similarity = threshold_similarity
        self.baseline_metrics: Optional[ResponseMetrics] = None

    def extract_metrics(self, status: int, headers: Dict, body: str, response_time_ms: int) -> ResponseMetrics:
        """
        Extract metrics from an HTTP response.
        
        Args:
            status: HTTP status code
            headers: Response headers dictionary
            body: Response body
            response_time_ms: Response time in milliseconds
            
        Returns:
            ResponseMetrics object
        """
        content_hash = hashlib.sha256(body.encode()).hexdigest()
        headers_str = json.dumps(headers, sort_keys=True)
        headers_hash = hashlib.sha256(headers_str.encode()).hexdigest()
        body_preview = body[:200] if body else ""

        return ResponseMetrics(
            status_code=status,
            content_hash=content_hash,
            content_length=len(body),
            response_time_ms=response_time_ms,
            headers_hash=headers_hash,
            body_preview=body_preview
        )

    def set_baseline(self, status: int, headers: Dict, body: str, response_time_ms: int) -> None:
        """
        Set the baseline response for comparison.
        
        Args:
            status: HTTP status code
            headers: Response headers
            body: Response body
            response_time_ms: Response time in milliseconds
        """
        self.baseline_metrics = self.extract_metrics(status, headers, body, response_time_ms)

    def analyze(self, status: int, headers: Dict, body: str, response_time_ms: int) -> Dict:
        """
        Analyze a response against the baseline.
        
        Args:
            status: HTTP status code
            headers: Response headers
            body: Response body
            response_time_ms: Response time in milliseconds
            
        Returns:
            Dictionary containing analysis results
        """
        if not self.baseline_metrics:
            raise ValueError("Baseline not set. Call set_baseline() first.")

        mutated_metrics = self.extract_metrics(status, headers, body, response_time_ms)

        analysis = {
            "status_changed": self.baseline_metrics.status_code != mutated_metrics.status_code,
            "content_hash_changed": self.baseline_metrics.content_hash != mutated_metrics.content_hash,
            "content_length_diff": mutated_metrics.content_length - self.baseline_metrics.content_length,
            "response_time_diff_ms": mutated_metrics.response_time_ms - self.baseline_metrics.response_time_ms,
            "headers_changed": self.baseline_metrics.headers_hash != mutated_metrics.headers_hash,
            "similarity_score": self._calculate_similarity(self.baseline_metrics.body_preview, mutated_metrics.body_preview),
            "body_diff": self._generate_diff(self.baseline_metrics.body_preview, mutated_metrics.body_preview),
            "vulnerabilities": self._detect_vulnerabilities(analysis={
                "status_changed": True,
                "content_hash_changed": mutated_metrics.content_hash != self.baseline_metrics.content_hash,
                "response_time_diff_ms": mutated_metrics.response_time_ms - self.baseline_metrics.response_time_ms,
            })
        }

        return analysis

    def _calculate_similarity(self, text1: str, text2: str) -> float:
        """
        Calculate similarity score between two texts (0-1).
        
        Args:
            text1: First text
            text2: Second text
            
        Returns:
            Similarity score
        """
        matcher = difflib.SequenceMatcher(None, text1, text2)
        return matcher.ratio()

    def _generate_diff(self, baseline: str, mutated: str) -> str:
        """
        Generate a unified diff between baseline and mutated responses.
        
        Args:
            baseline: Baseline response preview
            mutated: Mutated response preview
            
        Returns:
            Diff string
        """
        diff = list(difflib.unified_diff(
            baseline.splitlines(keepends=True),
            mutated.splitlines(keepends=True),
            lineterm=''
        ))
        return ''.join(diff[:20])  # Limit to first 20 lines

    def _detect_vulnerabilities(self, analysis: Dict) -> List[Dict]:
        """
        Detect potential vulnerabilities based on analysis.
        
        Args:
            analysis: Analysis results dictionary
            
        Returns:
            List of detected vulnerabilities
        """
        vulnerabilities = []

        # Timing-based attack detection
        if abs(analysis["response_time_diff_ms"]) > 1000:
            vulnerabilities.append({
                "type": VulnerabilityType.TIMING_ATTACK.value,
                "severity": "MEDIUM",
                "description": "Significant response time difference detected (possible timing attack)",
                "confidence": 0.7
            })

        # Content change detection
        if analysis["content_hash_changed"]:
            vulnerabilities.append({
                "type": VulnerabilityType.RESPONSE_DIFF.value,
                "severity": "MEDIUM",
                "description": "Response content differs from baseline",
                "confidence": 0.8
            })

        return vulnerabilities
