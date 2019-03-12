import uvm_pkg::*;
class test1 extends uvm_test;
`uvm_component_utils(test1)
function new(string name="test1", uvm_component parent);
    super.new(name, parent);
endfunction
virtual task main_phase(uvm_phase phase);
    phase.raise_objection(this);
    super.main_phase(phase);
    $display("run test1!");
    phase.drop_objection(this);
endtask
endclass