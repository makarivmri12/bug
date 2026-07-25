from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing import Optional, Dict, List
import logging

from .differential_analyzer import DifferentialAnalyzer, VulnerabilityType
from .advanced_state_tracker import AdvancedStateTracker
from .context_profiler import ContextProfiler

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="HLFA Heuristic Engine",
    description="Advanced heuristic analysis engine for vulnerability detection",
    version="1.0.0"
)

# Global instances
differential_analyzer = DifferentialAnalyzer()
state_tracker: Optional[AdvancedStateTracker] = None
context_profiler = ContextProfiler()


# Request/Response Models
class BaselineRequest(BaseModel):
    """Request to set baseline response"""
    status: int
    headers: Dict[str, str]
    body: str
    response_time_ms: int


class AnalysisRequest(BaseModel):
    """Request to analyze a response"""
    status: int
    headers: Dict[str, str]
    body: str
    response_time_ms: int
    mutation_type: Optional[str] = None


class StateTrackingRequest(BaseModel):
    """Request to track state"""
    user_id: str
    tokens: Optional[Dict[str, str]] = None
    cookies: Optional[Dict[str, str]] = None
    headers: Optional[Dict[str, str]] = None
    metadata: Optional[Dict] = None


class ProfileEndpointRequest(BaseModel):
    """Request to profile endpoint"""
    url: str
    method: str
    response_body: Optional[str] = ""
    response_headers: Optional[Dict[str, str]] = None


# Health Check Endpoint
@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "HLFA Heuristic Engine"}


# Differential Analysis Endpoints
@app.post("/api/v1/differential/baseline")
async def set_baseline(request: BaselineRequest):
    """
    Set the baseline response for differential analysis.
    
    Args:
        request: BaselineRequest containing response details
        
    Returns:
        Confirmation message
    """
    try:
        differential_analyzer.set_baseline(
            status=request.status,
            headers=request.headers,
            body=request.body,
            response_time_ms=request.response_time_ms
        )
        return {
            "status": "success",
            "message": "Baseline response set successfully"
        }
    except Exception as e:
        logger.error(f"Error setting baseline: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/differential/analyze")
async def analyze_response(request: AnalysisRequest):
    """
    Analyze a response against the baseline.
    
    Args:
        request: AnalysisRequest containing response details
        
    Returns:
        Analysis results with detected vulnerabilities
    """
    try:
        analysis = differential_analyzer.analyze(
            status=request.status,
            headers=request.headers,
            body=request.body,
            response_time_ms=request.response_time_ms
        )
        return {
            "status": "success",
            "analysis": analysis,
            "mutation_type": request.mutation_type
        }
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        logger.error(f"Error analyzing response: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


# State Tracking Endpoints
@app.post("/api/v1/state/capture")
async def capture_state(request: StateTrackingRequest):
    """
    Capture and track session state.
    
    Args:
        request: StateTrackingRequest containing state data
        
    Returns:
        State snapshot and detected changes
    """
    global state_tracker
    
    try:
        if not state_tracker or state_tracker.user_id != request.user_id:
            state_tracker = AdvancedStateTracker(request.user_id)
        
        snapshot = state_tracker.capture_state(
            tokens=request.tokens,
            cookies=request.cookies,
            headers=request.headers,
            metadata=request.metadata
        )
        
        return {
            "status": "success",
            "snapshot": snapshot.to_dict(),
            "changes": state_tracker.get_state_changes(),
            "anomalies": state_tracker.detect_anomalies()
        }
    except Exception as e:
        logger.error(f"Error capturing state: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/state/timeline/{user_id}")
async def get_state_timeline(user_id: str):
    """
    Get the state timeline for a user.
    
    Args:
        user_id: User/session identifier
        
    Returns:
        Complete state timeline
    """
    if not state_tracker or state_tracker.user_id != user_id:
        raise HTTPException(status_code=404, detail="No state timeline found")
    
    try:
        return {
            "status": "success",
            "user_id": user_id,
            "timeline": state_tracker.get_state_timeline(),
            "changes_count": len(state_tracker.changes)
        }
    except Exception as e:
        logger.error(f"Error retrieving state timeline: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


# Context Profiling Endpoints
@app.post("/api/v1/profile/endpoint")
async def profile_endpoint(request: ProfileEndpointRequest):
    """
    Profile an endpoint to extract context information.
    
    Args:
        request: ProfileEndpointRequest containing endpoint details
        
    Returns:
        Profiled endpoint context
    """
    try:
        context = context_profiler.profile_endpoint(
            url=request.url,
            method=request.method,
            response_body=request.response_body,
            response_headers=request.response_headers
        )
        
        return {
            "status": "success",
            "context": {
                "url": context.url,
                "method": context.method,
                "path": context.path,
                "parameters": dict(context.parameters),
                "authentication_required": context.authentication_required,
                "sensitive_fields": context.sensitive_fields,
                "input_types": context.input_types
            }
        }
    except Exception as e:
        logger.error(f"Error profiling endpoint: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/profile/endpoints")
async def get_all_contexts():
    """
    Get all profiled endpoint contexts.
    
    Returns:
        List of all profiled contexts
    """
    try:
        contexts = context_profiler.get_all_contexts()
        return {
            "status": "success",
            "total": len(contexts),
            "contexts": [
                {
                    "url": ctx.url,
                    "method": ctx.method,
                    "sensitive_fields": ctx.sensitive_fields,
                    "authentication_required": ctx.authentication_required
                }
                for ctx in contexts
            ]
        }
    except Exception as e:
        logger.error(f"Error retrieving contexts: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)
