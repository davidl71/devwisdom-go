#!/usr/bin/env python3
"""
Run all four exarp planning tools for devwisdom-go
"""

import json
import sys
from pathlib import Path

# Add parent project to path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

# Set working directory to devwisdom-go
devwisdom_path = Path(__file__).parent
import os
os.chdir(devwisdom_path.parent)  # Change to project root for exarp tools

from project_management_automation.tools.consolidated import (
    report,
    task_discovery,
    task_analysis,
    analyze_alignment
)
from project_management_automation.utils import find_project_root

def main():
    print("=" * 70)
    print("devwisdom-go - Exarp Planning Tools Execution")
    print("=" * 70)
    print()
    
    results = {}
    
    # 1. Generate PRD
    print("1️⃣  Generating PRD from PROJECT_GOALS.md...")
    try:
        # Generate PRD - note: generate_prd expects output_path as separate parameter
        from project_management_automation.tools.prd_generator import generate_prd
        prd_result = generate_prd(
            project_name="devwisdom-go",
            output_path=str(devwisdom_path / "PRD.md"),
            include_tasks=True,
            include_architecture=True,
            include_metrics=True
        )
        prd_data = json.loads(prd_result) if isinstance(prd_result, str) else prd_result
        results['prd'] = prd_data
        print(f"   ✅ PRD generated: {devwisdom_path / 'PRD.md'}")
    except Exception as e:
        print(f"   ❌ PRD generation failed: {e}")
        results['prd'] = {"error": str(e)}
    
    print()
    
    # 2. Discover tasks from TODO.md
    print("2️⃣  Discovering tasks from TODO.md...")
    try:
        discovery_result = task_discovery(
            action="markdown",
            doc_path=str(devwisdom_path / "TODO.md"),
            output_path=str(devwisdom_path / "discovered_tasks.json"),
            create_tasks=False  # Don't auto-create, just discover
        )
        discovery_data = json.loads(discovery_result) if isinstance(discovery_result, str) else discovery_result
        results['task_discovery'] = discovery_data
        total = discovery_data.get('summary', {}).get('total', 0)
        print(f"   ✅ Discovered {total} tasks")
    except Exception as e:
        print(f"   ❌ Task discovery failed: {e}")
        results['task_discovery'] = {"error": str(e)}
    
    print()
    
    # 3. Analyze task structure
    print("3️⃣  Analyzing task structure (hierarchy, tags, duplicates)...")
    try:
        # Analyze hierarchy - import directly to avoid type issues
        from project_management_automation.tools.task_hierarchy_analyzer import analyze_task_hierarchy
        from typing import Optional
        hierarchy_result = analyze_task_hierarchy(
            output_format="json",
            output_path=str(devwisdom_path / "task_hierarchy_analysis.json"),
            include_recommendations=True
        )
        hierarchy_data = hierarchy_result if isinstance(hierarchy_result, dict) else json.loads(hierarchy_result)
        
        # Save hierarchy results
        with open(devwisdom_path / "task_hierarchy_analysis.json", 'w') as f:
            json.dump(hierarchy_data, f, indent=2)
        
        # Analyze duplicates
        duplicates_result = task_analysis(
            action="duplicates",
            similarity_threshold=0.85,
            auto_fix=False,
            output_path=str(devwisdom_path / "task_duplicates_analysis.json")
        )
        duplicates_data = json.loads(duplicates_result) if isinstance(duplicates_result, str) else duplicates_result
        
        # Analyze tags
        tags_result = task_analysis(
            action="tags",
            dry_run=True,
            output_path=str(devwisdom_path / "task_tags_analysis.json")
        )
        tags_data = json.loads(tags_result) if isinstance(tags_result, str) else tags_result
        
        results['task_analysis'] = {
            'hierarchy': hierarchy_data,
            'duplicates': duplicates_data,
            'tags': tags_data
        }
        print(f"   ✅ Task analysis complete")
    except Exception as e:
        print(f"   ❌ Task analysis failed: {e}")
        results['task_analysis'] = {"error": str(e)}
    
    print()
    
    # 4. Check alignment with PROJECT_GOALS.md
    print("4️⃣  Checking alignment with PROJECT_GOALS.md...")
    try:
        # Temporarily copy PROJECT_GOALS.md to project root for alignment tool
        # (It looks for it in project root)
        goals_source = devwisdom_path / "PROJECT_GOALS.md"
        goals_target = project_root / "PROJECT_GOALS.md"
        
        # Backup original if exists
        original_exists = goals_target.exists()
        original_backup = None
        if original_exists:
            original_backup = project_root / "PROJECT_GOALS.md.backup"
            import shutil
            shutil.copy2(goals_target, original_backup)
        
        # Copy devwisdom-go goals
        import shutil
        shutil.copy2(goals_source, goals_target)
        
        try:
            alignment_result = analyze_alignment(
                action="todo2",
                create_followup_tasks=False,  # Don't create tasks automatically
                output_path=str(devwisdom_path / "alignment_report.md")
            )
            alignment_data = json.loads(alignment_result) if isinstance(alignment_result, str) else alignment_result
            results['alignment'] = alignment_data
            
            score = alignment_data.get('data', {}).get('average_alignment_score', 0)
            misaligned = alignment_data.get('data', {}).get('misaligned_count', 0)
            print(f"   ✅ Alignment check complete")
            print(f"      Alignment Score: {score:.1f}%")
            print(f"      Misaligned Tasks: {misaligned}")
        finally:
            # Restore original or remove copy
            if original_backup and original_backup.exists():
                shutil.move(original_backup, goals_target)
            elif goals_target.exists():
                goals_target.unlink()
                
    except Exception as e:
        print(f"   ❌ Alignment check failed: {e}")
        import traceback
        traceback.print_exc()
        results['alignment'] = {"error": str(e)}
    
    print()
    print("=" * 70)
    print("Summary")
    print("=" * 70)
    print()
    
    # Save combined results
    output_file = devwisdom_path / "exarp_planning_results.json"
    with open(output_file, 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"📊 Results saved to: {output_file}")
    print()
    print("Generated Files:")
    print(f"  - {devwisdom_path / 'PRD.md'}")
    print(f"  - {devwisdom_path / 'discovered_tasks.json'}")
    print(f"  - {devwisdom_path / 'task_hierarchy_analysis.json'}")
    print(f"  - {devwisdom_path / 'task_duplicates_analysis.json'}")
    print(f"  - {devwisdom_path / 'task_tags_analysis.json'}")
    print(f"  - {devwisdom_path / 'alignment_report.md'}")
    print(f"  - {output_file}")
    print()
    print("✅ All planning tools executed!")

if __name__ == "__main__":
    main()
