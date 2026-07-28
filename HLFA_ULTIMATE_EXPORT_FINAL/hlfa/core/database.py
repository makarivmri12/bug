from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker, declarative_base
from sqlalchemy import Column, String, Float, DateTime, Text
from datetime import datetime
import os

Base = declarative_base()
DB_URL = "sqlite+aiosqlite:///hlfa_data.db"

class FindingModel(Base):
    __tablename__ = "findings"
    id = Column(String, primary_key=True)
    target = Column(String)
    category = Column(String)
    severity = Column(String)
    cvss_score = Column(Float)
    evidence_path = Column(String)
    timestamp = Column(DateTime, default=datetime.utcnow)

engine = create_async_engine(DB_URL, echo=False)
async_session = sessionmaker(engine, class_=AsyncSession, expire_on_commit=False)

async def init_db():
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
