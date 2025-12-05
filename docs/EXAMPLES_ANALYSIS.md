# Example Expansion Opportunities

This document summarizes what we can build next on top of the existing example set that lives in [`/examples`](../examples) and mirrors the fork at <https://github.com/AndreyLeonov80/sdfx/tree/master/examples>.

## 1. Repository Inventory

* We cloned Andrey Leonov's fork under `/tmp/analyze_sdfx` and compared it with the current branch (`feat/analyze-create-sdfx-examples`). `diff -qr /home/engine/project/examples /tmp/analyze_sdfx/examples` produced no output, which confirms both trees are identical.
* Running `find /home/engine/project/examples -mindepth 1 -maxdepth 1 -type d | wc -l` shows **73** concrete example directories. They already cover several recurring themes:
  * **Mechanical hardware** – e.g. `3dp_nutbolt`, `gears`, `nutsandbolts`, `sprue`.
  * **Electronics & IoT enclosures** – `eurorack`, `maestro`, `maixgo`, `panel_box`, `pico_cnc`.
  * **Grid / storage tooling** – `gridfinity`, `bolt_container`, `tabbox`.
  * **Mathematical & artistic solids** – `gyroid`, `bucky`, `voronoi`, `bezier`, `spiral`.
  * **2D DXF/SVG work** – `hole_patterns`, `spiral`, `text`, `dc2test`.
  * **Mesh post-processing** – `hollowing_stl`, `gyroid`’s teapot demo, `monkey_hat`.

## 2. Observations About Current Coverage

### 2.1 Under-exercised building blocks

* `obj.SpringParms` (`obj/spring.go`) only appears once in the tree (`examples/pico_cnc/penholder.go`). There is no focused example that tunes wall thickness, boss lengths, or demonstrates both `Spring2D` and `Spring3D` outputs.
* `obj.ImportSTL` is available (`obj/stl.go`) but only used in three places: `examples/gyroid/main.go`, `examples/hollowing_stl/main.go`, and `examples/monkey_hat/main.go`. None of them combines multiple imported meshes or mixes lattice infill with external geometry.
* `render.ToTriangles` (used for in-memory mesh pipelines) shows up once in `examples/gyroid/main.go`. Users who want to post-process meshes in Go do not have a clean reference.

### 2.2 Output formats and renderers

* `render.To3MF` is only invoked in `examples/bucky/main.go` (lines 75–82). We are not showcasing 3MF’s benefits (multi-part builds, metadata) elsewhere.
* `render.NewDualContour(...)` is referenced by just three examples (`dc2test`, `pool`, `monkey_hat`). Most models stick to marching cubes, so users do not see when dual contouring produces cleaner surfaces or fewer triangles.
* DXF/SVG generation exists (`render.ToDXF`, `render.ToSVG` in `render/render.go`), but only seven examples output DXF and none emit SVG by default.

### 2.3 Parametric workflows missing from examples

* Storage-focused examples (`gridfinity`, `bolt_container`) do not yet demonstrate advanced variants such as tilted drawers, built-in labels, or magnet/fastener placement matrices.
* Enclosure examples lean on `obj.Standoff`/`obj.PanelBox` primitives but do not include flexible closure mechanisms (living hinges, snap latches) even though `obj.Spring` makes those attainable.
* There is no side-by-side benchmark example that logs triangle counts / timing for different renderers or resolution settings, which would help users tune slicing fidelity.

## 3. Candidate Example Backlog

| Working name | Based on (existing dirs) | Idea summary | Key APIs to exercise | Suggested outputs |
| --- | --- | --- | --- | --- |
| **`gridfinity_drawer`** | `examples/gridfinity` | Extend the existing base/body demo with stackable drawers: add tilted fronts, dividers, and label plates controlled by parameters. | `obj.GfBase`, `obj.GfBody`, boolean ops from `sdf`, simple Bézier lofts for the chamfered front. | Multiple STLs (`drawer_shell.stl`, `drawer_divider.stl`). |
| **`flex_latch_kit`** | `examples/pico_cnc`, `examples/panel_box` | Generate a small library of printable snap latches and living hinges that enclosure authors can drop into their designs. Demonstrate both 2D profiles (`Spring2D`) and 3D extrusions (`Spring3D`). | `obj.SpringParms`, `sdf.Extrude3D`, `obj.Standoff`, filleting helpers (`sdf.Offset3D`). | STL set for latch variants + DXF of the 2D profiles. |
| **`hardware_pack_3mf`** | `examples/3dp_nutbolt`, `examples/bucky` | Produce a color-coded fastener kit exported as a single 3MF with multiple components (nuts, bolts, washers). Show how to pack items as separate objects inside one 3MF scene. | `obj.HexHead`, `obj.Nut`, `render.To3MF`, go3mf metadata assignment. | Single `fasteners.3mf` plus optional STL fallback. |
| **`surface_compare`** | `examples/dc2test`, `examples/monkey_hat` | Render the same complex model twice (marching cubes octree vs dual contour) and print triangle counts + timings. Helps users pick algorithms per shape. | `render.NewMarchingCubesOctree`, `render.NewDualContour`, `render.ToTriangles`, metrics from `render.Render3.Info`. | STL pairs and a markdown/CSV report. |
| **`mesh_ops_pipeline`** | `examples/hollowing_stl`, `examples/gyroid` | Import multiple STLs (e.g., `files/monkey.stl`, `files/teapot.stl`), apply boolean trims, punch registration sockets, and re-export. Demonstrates `obj.ImportSTL` as part of a richer workflow. | `obj.ImportSTL`, `sdf.Union3D`/`Difference3D`, `sdf.Shell3D`, `render.ToSTL`. | Intermediate/ final STLs for each stage. |

## 4. Next Steps

1. Prioritize the ideas above based on immediate needs (e.g., marketing wants Gridfinity content, or we need dual-contour guidance for support tickets).
2. For each selected idea, create a new directory under `/examples`, reuse the conventions already present (Makefile + SHA1SUM + `main.go`).
3. When applicable, add small README snippets or inline comments explaining novel parameters (for instance, why a specific `SpringParms` combination gives the correct clip stiffness).
4. Capture rendered artifacts (PNG screenshots or STL previews) in the `docs/gallery` folder to advertise the new examples once they land.

These additions would keep parity with Andrey Leonov’s fork while also showcasing parts of the API surface that are currently invisible to users.
