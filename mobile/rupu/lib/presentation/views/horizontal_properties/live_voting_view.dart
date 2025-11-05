import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rupu/presentation/views/horizontal_properties/controllers/live_voting_controller.dart';

class LiveVotingView extends StatelessWidget {
  final int propertyId;
  final int votingGroupId;
  final int votingId;

  const LiveVotingView({
    super.key,
    required this.propertyId,
    required this.votingGroupId,
    required this.votingId,
  });

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<LiveVotingController>();
    return Scaffold(
      appBar: AppBar(
        title: const Text('Votación en vivo'),
      ),
      body: const Center(
        child: Text('Live Voting View'),
      ),
    );
  }
}
