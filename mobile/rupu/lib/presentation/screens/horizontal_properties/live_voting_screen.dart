import 'package:flutter/material.dart';
import 'package:rupu/config/routers/app_bindings.dart';
import 'package:rupu/presentation/views/horizontal_properties/live_voting_view.dart';

class LiveVotingScreen extends StatelessWidget {
  static const name = 'live-voting-screen';
  final int pageIndex;
  final int propertyId;
  final int votingGroupId;
  final int votingId;

  const LiveVotingScreen({
    super.key,
    required this.pageIndex,
    required this.propertyId,
    required this.votingGroupId,
    required this.votingId,
  });

  @override
  Widget build(BuildContext context) {
    LiveVotingBinding.register(
      propertyId: propertyId,
      votingGroupId: votingGroupId,
      votingId: votingId,
    );
    return LiveVotingView(
      propertyId: propertyId,
      votingGroupId: votingGroupId,
      votingId: votingId,
    );
  }
}
