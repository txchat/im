local name = "im"
local p_im = Proto(name, "im layer protocal");

local f_pack_size = ProtoField.uint32(name .. ".packsize", "Packet Size", base.DEC)
local f_header_len = ProtoField.uint16(name .. ".header", "Header Length", base.DEC)
local f_version = ProtoField.uint16(name .. ".version", "Version", base.DEC)
local f_option = ProtoField.uint32(name .. ".option", "Option", base.DEC)
local f_seq = ProtoField.uint32(name .. ".seq", "Seq", base.DEC)
local f_ack = ProtoField.uint32(name .. ".ack", "Ack", base.DEC)
local f_data = ProtoField.bytes(name .. ".data", "Payload", FT_BYTES)

p_im.fields = { f_pack_size,f_header_len,f_version,f_option,f_seq,f_ack,f_data }

local option_name = {
        [0] = 'Undefined',
        [1] = 'Auth',
        [2] = 'AuthReply',
        [3] = 'Heartbeat',
        [4] = 'HeartbeatReply',
        [5] = 'Disconnect',
        [6] = 'DisconnectReply',
        [7] = 'SendMsg',
        [8] = 'SendMsgReply',
        [9] = 'ReceiveMsg',
        [10] = 'ReceiveMsgReply',
        [14] = 'SyncMsgReq',
        [15] = 'SyncMsgReply'
    }

local dtalk_proto = Dissector.get("dtalk")
local data_dis = Dissector.get("data")

function p_im.dissector(buf, pkt, tree)
        local offset = 0
        local subtree = tree:add(p_im, buf())
        pkt.cols.protocol = p_im.name

        --packate size
        local pkglen_tvbr = buf(offset, 4)
        local pkglen = pkglen_tvbr:uint()
        subtree:add(f_pack_size, pkglen_tvbr)
        offset = offset + 4

        if pkglen ~= buf:len() then
                data_dis:call(buf(2):tvb(), pkt, tree)
        else
                --deal im parse
                -- header Length
                subtree:add(f_header_len, buf(offset, 2))
                offset = offset + 2
                --version
                subtree:add(f_version, buf(offset, 2))
                offset = offset + 2
                --option
                local opt_tvbr = buf(offset, 4)
                local optcode = opt_tvbr:uint()
                local opttree = subtree:add(f_option, opt_tvbr)
                offset = offset + 4
                opttree:append_text(' (' .. option_name[optcode] .. ') ')
                pkt.private["dtalk_opt_type"] = optcode
                --sequence
                subtree:add(f_seq, buf(offset, 4))
                offset = offset + 4
                --acknowledge
                subtree:add(f_ack, buf(offset, 4))
                offset = offset + 4
                --data
                if buf:len()-offset > 0 then
                        subtree:add(f_data, buf(offset, buf:len()-offset))
                        --sub protocal
                        if dtalk_proto ~= nil then
                                dtalk_proto:call(buf(offset):tvb(), pkt, tree)
                        end
                else
                        subtree:add('Payload:', 'empty')
                end  
        end
end

local websocket_table = DissectorTable.get("ws.port")

websocket_table:add(3102, p_im)