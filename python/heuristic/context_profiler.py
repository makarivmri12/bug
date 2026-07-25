from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass
import json
from urllib.parse import urlparse, parse_qs
import re


@dataclass
class EndpointContext:
    """Context information for an endpoint"""
    url: str
    method: str
    path: str
    parameters: Dict[str, List[str]]
    authentication_required: bool
    sensitive_fields: List[str]
    input_types: Dict[str, str]


class ContextProfiler:
    """
    Analyzes endpoint context to identify sensitive areas and attack vectors.
    """

    # Patterns for sensitive data
    SENSITIVE_PATTERNS = {
        'id_field': re.compile(r'(user_?id|account_?id|customer_?id|id)$', re.IGNORECASE),
        'auth_field': re.compile(r'(token|api_?key|secret|password|auth)', re.IGNORECASE),
        'email_field': re.compile(r'(email|e_?mail|contact)', re.IGNORECASE),
        'amount_field': re.compile(r'(amount|price|cost|total|value)', re.IGNORECASE),
    }

    def __init__(self):
        self.endpoint_contexts: Dict[str, EndpointContext] = {}

    def profile_endpoint(self, url: str, method: str, response_body: str = "", response_headers: Dict = None) -> EndpointContext:
        """
        Profile an endpoint to extract context information.
        
        Args:
            url: Endpoint URL
            method: HTTP method
            response_body: Response body content
            response_headers: Response headers
            
        Returns:
            EndpointContext
        """
        parsed_url = urlparse(url)
        path = parsed_url.path
        parameters = parse_qs(parsed_url.query)

        # Extract parameters from path
        path_params = self._extract_path_parameters(path)
        parameters.update(path_params)

        # Detect authentication requirements
        auth_required = self._detect_auth_requirement(response_body, response_headers or {})

        # Identify sensitive fields
        sensitive_fields = self._identify_sensitive_fields(list(parameters.keys()))

        # Detect input types
        input_types = self._detect_input_types(parameters)

        context = EndpointContext(
            url=url,
            method=method,
            path=path,
            parameters=parameters,
            authentication_required=auth_required,
            sensitive_fields=sensitive_fields,
            input_types=input_types
        )

        self.endpoint_contexts[url] = context
        return context

    def _extract_path_parameters(self, path: str) -> Dict[str, List[str]]:
        """
        Extract parameters from URL path (e.g., /users/{id}/posts/{post_id}).
        
        Args:
            path: URL path
            
        Returns:
            Dictionary of parameters
        """
        params = {}
        path_segments = path.split('/')
        
        for segment in path_segments:
            if segment.startswith('{') and segment.endswith('}'):
                param_name = segment[1:-1]
                params[param_name] = ["<path_parameter>"]
            elif segment.isdigit() or segment.lower() in ['new', 'edit', 'delete']:
                params[f"path_segment_{segment}"] = [segment]
        
        return params

    def _detect_auth_requirement(self, response_body: str, headers: Dict) -> bool:
        """
        Detect if endpoint requires authentication.
        
        Args:
            response_body: Response body
            headers: Response headers
            
        Returns:
            True if authentication appears required
        """
        # Check for common auth indicators
        auth_indicators = [
            'unauthorized',
            'forbidden',
            'authentication required',
            'login',
            'please sign in',
            'www-authenticate'
        ]
        
        response_text = (response_body + json.dumps(headers)).lower()
        return any(indicator in response_text for indicator in auth_indicators)

    def _identify_sensitive_fields(self, fields: List[str]) -> List[str]:
        """
        Identify sensitive fields based on name patterns.
        
        Args:
            fields: List of field names
            
        Returns:
            List of identified sensitive fields
        """
        sensitive = []
        
        for field in fields:
            for pattern_name, pattern in self.SENSITIVE_PATTERNS.items():
                if pattern.search(field):
                    sensitive.append(field)
                    break
        
        return sensitive

    def _detect_input_types(self, parameters: Dict[str, List[str]]) -> Dict[str, str]:
        """
        Detect input types for parameters.
        
        Args:
            parameters: Parameters dictionary
            
        Returns:
            Dictionary mapping parameter names to detected types
        """
        input_types = {}
        
        for param_name, values in parameters.items():
            if values and len(values) > 0:
                value = values[0]
                
                if value.isdigit():
                    input_types[param_name] = "integer"
                elif value.lower() in ['true', 'false']:
                    input_types[param_name] = "boolean"
                elif '@' in value:
                    input_types[param_name] = "email"
                else:
                    input_types[param_name] = "string"
            else:
                input_types[param_name] = "unknown"
        
        return input_types

    def get_context(self, url: str) -> Optional[EndpointContext]:
        """
        Get profiled context for an endpoint.
        
        Args:
            url: Endpoint URL
            
        Returns:
            EndpointContext or None
        """
        return self.endpoint_contexts.get(url)

    def get_all_contexts(self) -> List[EndpointContext]:
        """
        Get all profiled endpoint contexts.
        
        Returns:
            List of EndpointContext objects
        """
        return list(self.endpoint_contexts.values())
