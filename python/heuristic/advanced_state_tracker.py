from typing import Dict, List, Optional, Any
from dataclasses import dataclass, field
from enum import Enum
import time
import json


class StateChangeType(Enum):
    """Types of state changes"""
    TOKEN_ADDED = "token_added"
    TOKEN_MODIFIED = "token_modified"
    COOKIE_ADDED = "cookie_added"
    COOKIE_MODIFIED = "cookie_modified"
    HEADER_ADDED = "header_added"
    HEADER_MODIFIED = "header_modified"
    METADATA_ADDED = "metadata_added"


@dataclass
class StateSnapshot:
    """Snapshot of state at a point in time"""
    timestamp: float
    tokens: Dict[str, str] = field(default_factory=dict)
    cookies: Dict[str, str] = field(default_factory=dict)
    headers: Dict[str, str] = field(default_factory=dict)
    metadata: Dict[str, Any] = field(default_factory=dict)

    def to_dict(self) -> Dict:
        """Convert snapshot to dictionary"""
        return {
            "timestamp": self.timestamp,
            "tokens": self.tokens,
            "cookies": self.cookies,
            "headers": self.headers,
            "metadata": self.metadata
        }


@dataclass
class StateChange:
    """Represents a single state change"""
    change_type: StateChangeType
    key: str
    old_value: Optional[str]
    new_value: str
    timestamp: float


class AdvancedStateTracker:
    """
    Advanced session state tracker that monitors and records state changes.
    Detects state mutations, correlations, and anomalies.
    """

    def __init__(self, user_id: str):
        """
        Initialize the state tracker.
        
        Args:
            user_id: Unique user/session identifier
        """
        self.user_id = user_id
        self.snapshots: List[StateSnapshot] = []
        self.changes: List[StateChange] = []
        self.current_state = StateSnapshot(timestamp=time.time())

    def capture_state(self, tokens: Dict = None, cookies: Dict = None, headers: Dict = None, metadata: Dict = None) -> StateSnapshot:
        """
        Capture the current state and detect changes.
        
        Args:
            tokens: Authentication tokens
            cookies: HTTP cookies
            headers: HTTP headers
            metadata: Arbitrary metadata
            
        Returns:
            StateSnapshot
        """
        new_state = StateSnapshot(
            timestamp=time.time(),
            tokens=tokens or self.current_state.tokens.copy(),
            cookies=cookies or self.current_state.cookies.copy(),
            headers=headers or self.current_state.headers.copy(),
            metadata=metadata or self.current_state.metadata.copy()
        )

        # Detect changes
        self._detect_changes(self.current_state, new_state)

        self.snapshots.append(new_state)
        self.current_state = new_state

        return new_state

    def _detect_changes(self, old: StateSnapshot, new: StateSnapshot) -> None:
        """
        Detect and record state changes.
        
        Args:
            old: Previous state snapshot
            new: New state snapshot
        """
        # Check token changes
        for key, value in new.tokens.items():
            if key not in old.tokens:
                self.changes.append(StateChange(
                    change_type=StateChangeType.TOKEN_ADDED,
                    key=key,
                    old_value=None,
                    new_value=value,
                    timestamp=new.timestamp
                ))
            elif old.tokens[key] != value:
                self.changes.append(StateChange(
                    change_type=StateChangeType.TOKEN_MODIFIED,
                    key=key,
                    old_value=old.tokens[key],
                    new_value=value,
                    timestamp=new.timestamp
                ))

        # Check cookie changes
        for key, value in new.cookies.items():
            if key not in old.cookies:
                self.changes.append(StateChange(
                    change_type=StateChangeType.COOKIE_ADDED,
                    key=key,
                    old_value=None,
                    new_value=value,
                    timestamp=new.timestamp
                ))
            elif old.cookies[key] != value:
                self.changes.append(StateChange(
                    change_type=StateChangeType.COOKIE_MODIFIED,
                    key=key,
                    old_value=old.cookies[key],
                    new_value=value,
                    timestamp=new.timestamp
                ))

    def get_state_timeline(self) -> List[Dict]:
        """
        Get the complete state timeline.
        
        Returns:
            List of state snapshots as dictionaries
        """
        return [snapshot.to_dict() for snapshot in self.snapshots]

    def get_state_changes(self) -> List[Dict]:
        """
        Get all recorded state changes.
        
        Returns:
            List of state changes
        """
        return [
            {
                "type": change.change_type.value,
                "key": change.key,
                "old_value": change.old_value,
                "new_value": change.new_value,
                "timestamp": change.timestamp
            }
            for change in self.changes
        ]

    def detect_anomalies(self) -> List[Dict]:
        """
        Detect anomalies in state changes.
        
        Returns:
            List of detected anomalies
        """
        anomalies = []

        # Check for too many rapid changes
        if len(self.changes) > 50:
            anomalies.append({
                "type": "HIGH_CHANGE_RATE",
                "description": f"Unusually high number of state changes: {len(self.changes)}",
                "severity": "MEDIUM"
            })

        # Check for token/cookie proliferation
        if len(self.current_state.tokens) > 10:
            anomalies.append({
                "type": "TOKEN_PROLIFERATION",
                "description": f"Unusually many tokens: {len(self.current_state.tokens)}",
                "severity": "LOW"
            })

        return anomalies
