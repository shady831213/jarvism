import uvm_pkg::*;
class test1 extends uvm_test;
`uvm_component_utils(test1)
function new(string name="test1", uvm_component parent);
    super.new(name, parent);
endfunction
virtual task main_phase(uvm_phase phase);
    phase.raise_objection(this);
    super.main_phase(phase);
    randcase
        1: `uvm_error(this.get_name(), "test1")
        1: `uvm_fatal(this.get_name(), "test1")
        1: `uvm_warning(this.get_name(), "test1")
        1: `uvm_info(this.get_name(), "test1", UVM_LOW)
    endcase
    phase.drop_objection(this);
endtask
endclass